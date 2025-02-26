# linux 内存管理

每个项目的物理地址对于进程不可见，谁也不能直接访问这个物理地址。操作系统会给进程分配一个虚拟地址。所有进程看到的这个地址都是一样的，里面的内存都是从 0 开始编号。  
内存中的地址可分为 ： 
  逻辑地址：逻辑地址是由段和段偏移量组成。  
  线性地址：线性地址也称为虚拟地址，由十六进制表示，经过分页单元转化为物理之地。  
  物理地址：真是的物理地址。 
虚拟空间，一部分用来放内核的东西，称为内核空间，一部分用来放进程的东西，称为用户空间，低地址是用户态信息，高地址是内核信息。
## 逻辑地址映射到线性地址
![image](https://github.com/yanlongLv/notes/blob/main/linux/picture/duan1.jpg)   
逻辑地址由两部分组成，段标识符合和段内地址偏移量。 
分段机制下的虚拟地址由两部分组成，段选择子和段内偏移量。段选择子里面最重要的是段号，用作段表的索引。段表里面保存的是这个段的基地址、段的界限和特权等级等。虚拟地址中的段内偏移量应该位于 0 和段界限之间。  

## 虚拟地址到物理地址的映射    

从虚拟地址到物理地址的转换方式，称为分页（Paging）。  
![image](https://github.com/yanlongLv/notes/blob/main/linux/picture/xuni1.jpg)  
虚拟地址分为两部分，页号和页内偏移。页号作为页表的索引，页表包含物理页每页所在物理内存的基地址。这个基地址与页内偏移的组合就形成了物理内存地址。  
![image](https://github.com/yanlongLv/notes/blob/main/linux/picture/yishe.jpg)  
我们可以试着将页表再分页，4G 的空间需要 4M 的页表来存储映射。我们把这 4M 分成 1K（1024）个 4K，每个 4K 又能放在一页里面，这样 1K 个 4K 就是 1K 个页，这 1K 个页也需要一个表进行管理，我们称为页目录表，这个页目录表里面有 1K 项，每项 4 个字节，页目录表大小也是 4K。  
页目录有 1K 项，用 10 位就可以表示访问页目录的哪一项。这一项其实对应的是一整页的页表项，也即 4K 的页表项。每个页表项也是 4 个字节，因而一整页的页表项是 1K 个。再用 10 位就可以表示访问页表项的哪一项，页表项中的一项对应的就是一个页，是存放数据的页，这个页的大小是 4K，用 12 位可以定位这个页内的任何一个位置。  
![image](https://github.com/yanlongLv/notes/blob/main/linux/picture/page.jpg)  
如果只使用页表，也需要完整的 1M 个页表项共 4M 的内存，但是如果使用了页目录，页目录需要 1K 个全部分配，占用内存 4K，但是里面只有一项使用了。到了页表项，只需要分配能够管理那个数据页的页表项页就可以了，也就是说，最多 4K，这样内存就节省多了。  

如果要申请小块内存，就用 brk。如果申请一大块内存，就要用 mmap。对于堆的申请来讲，mmap 是映射内存空间到物理内存。

如果 PTE，也就是页表项，从来没有出现过，那就是新映射的页。如果是匿名页，就是第一种情况，应该映射到一个物理内存页，在这里调用的是 do_anonymous_page。如果是映射到文件，调用的就是 do_fault，这是第二种情况。如果 PTE 原来出现过，说明原来页面在物理内存中，后来换出到硬盘了，现在应该换回来，调用的是 do_swap_page。  
# 文件系统
硬盘分成相同大小的单元，我们称为块（Block）。一块的大小是扇区大小的整数倍，默认是 4K。另外，文件还有元数据部分，例如名字、权限等，这就需要一个结构 inode 来存放。
## inode介绍
```
struct ext4_inode {
  __le16  i_mode;    /* File mode */
  __le16  i_uid;    /* Low 16 bits of Owner Uid */
  __le32  i_size_lo;  /* Size in bytes */
  __le32  i_atime;  /* Access time */
  __le32  i_ctime;  /* Inode Change time */
  __le32  i_mtime;  /* Modification time */
  __le32  i_dtime;  /* Deletion Time */
  __le16  i_gid;    /* Low 16 bits of Group Id */
  __le16  i_links_count;  /* Links count */
  __le32  i_blocks_lo;  /* Blocks count */
  __le32  i_flags;  /* File flags */
......
  __le32  i_block[EXT4_N_BLOCKS];/* Pointers to blocks */
  __le32  i_generation;  /* File version (for NFS) */
  __le32  i_file_acl_lo;  /* File ACL */
  __le32  i_size_high;
......
};
```

node 里面有文件的读写权限 i_mode，属于哪个用户 i_uid，哪个组 i_gid，大小是多少 i_size_io，占用多少个块 i_blocks_io以及权限、用户、大小这些信息。i_block 表示某个文件分成几块、每一块在哪里。  
![image](https://github.com/yanlongLv/notes/blob/main/linux/picture/fileblock.jpg)  
如果一个文件比较大，12 块放不下。当我们用到 i_block[12]的时候，就不能直接放数据块的位置了，要不然 i_block 很快就会用完了。这该怎么办呢？我们需要想个办法。我们可以让 i_block[12]指向一个块，这个块里面不放数据块，而是放数据块的位置，这个块我们称为间接块。也就是说，我们在 i_block[12]里面放间接块的位置，通过 i_block[12]找到间接块后，间接块里面放数据块的位置，通过间接块可以找到数据块。如果文件再大一些，i_block[13]会指向一个块，我们可以用二次间接块。二次间接块里面存放了间接块的位置，间接块里面存放了数据块的位置，数据块里面存放的是真正的数据。如果文件再大一些，i_block[14]会指向三次间接块。原理和上面都是一样的，就像一层套一层的俄罗斯套娃，一层一层打开，才能拿到最中心的数据块。  
对于大文件来讲，我们要多次读取硬盘才能找到相应的块，这样访问速度就会比较慢。为了解决这个问题，ext4 做了一定的改变。它引入了一个新的概念，叫做 Extents。  
一个文件大小为 128M，如果使用 4k 大小的块进行存储，需要 32k 个块。如果按照 ext2 或者 ext3 那样散着放，数量太大了。但是 Extents 可以用于存放连续的块，也就是说，我们可以把 128M 放在一个 Extents 里面。这样的话，对大文件的读写性能提高了，文件碎片也减少了。  
![image](https://github.com/yanlongLv/notes/blob/main/linux/picture/exents.jpg)  
eh_entries 表示这个节点里面有多少项。这里的项分两种，如果是叶子节点，这一项会直接指向硬盘上的连续块的地址，我们称为数据节点 ext4_extent；如果是分支节点，这一项会指向下一层的分支节点或者叶子节点，我们称为索引节点 ext4_extent_idx。这两种类型的项的大小都是 12 个 byte。  
除了根节点，其他的节点都保存在一个块 4k 里面，4k 扣除 ext4_extent_header 的 12 个 byte，剩下的能够放 340 项，每个 extent 最大能表示 128MB 的数据，340 个 extent 会使你表示的文件达到 42.5GB。这已经非常大了，如果再大，我们可以增加树的深度。  
如果我要保存一个数据块，或者要保存一个 inode，我应该放在硬盘上的哪个位置呢？难道需要将所有的 inode 列表和块列表扫描一遍，找个空的地方随便放吗？当然，这样效率太低了。所以在文件系统里面，我们专门弄了一个块来保存 inode 的位图。在这 4k 里面，每一位对应一个 inode。如果是 1，表示这个 inode 已经被用了；如果是 0，则表示没被用。同样，我们也弄了一个块保存 block 的位图

# 网络系统  
select/poll/epoll都是IO多路复用机制，可以同时监控多个描述符，当某个描述符就绪(读或写就绪)，则立刻通知相应程序进行读或写操作。本质上select/poll/epoll都是同步I/O，即读写是阻塞的。  
## select
```
int select (int maxfd, 
            fd_set *readfds, 
            fd_set *writefds, 
            fd_set *exceptfds, 
            struct timeval *timeout);
```
maxfd：代表要监控的最大文件描述符fd+1  
writefds：监控可写fd  
readfds：监控可读fd  
exceptfds：监控异常fd  
timeout：超时时长  
   1. NULL，代表没有设置超时，则会一直阻塞直到文件描述符上的事件触发
   2. 0，代表不等待，立即返回，用于检测文件描述符状态
   3. 正整数，代表当指定时间没有事件触发，则超时返回  

select函数监控3类文件描述符，调用select函数后会阻塞，直到描述符fd准备就绪（有数据可读、可写、异常）或者超时，函数便返回。 当select函数返回后，可通过遍历描述符集合，找到就绪的描述符。

### select缺点
 1. 文件描述符个数受限：单进程能够监控的文件描述符的数量存在最大限制，在Linux上一般为1024，可以通过修改宏定义增大上限，但同样存在效率低的弱势;
 2. 性能衰减严重：IO随着监控的描述符数量增长，其性能会线性下降; 

 ## poll

 ```
 int poll (struct pollfd *fds, unsigned int nfds, int timeout);
 ```
其中pollfd表示监视的描述符集合，如下:
```
struct pollfd {
    int fd; //文件描述符
    short events; //监视的请求事件
    short revents; //已发生的事件
};
```
pollfd结构包含了要监视的event和发生的event，并且pollfd并没有最大数量限制。 和select函数一样，当poll函数返回后，可以通过遍历描述符集合，找到就绪的描述符。  
### poll缺点  
从上面看select和poll都需要在返回后，通过遍历文件描述符来获取已经就绪的socket。同时连接的大量客户端在同一时刻可能只有很少的处于就绪状态，因此随着监视的描述符数量的增长，其性能会线性下降。

## epoll
epoll是在内核2.6中提出的，是select和poll的增强版。相对于select和poll来说，epoll更加灵活，没有描述符数量限制。epoll使用一个文件描述符管理多个描述符，将用户空间的文件描述符的事件存放到内核的一个事件表中，这样在用户空间和内核空间的copy只需一次。epoll机制是Linux最高效的I/O复用机制，在一处等待多个文件句柄的I/O事件。

select/poll都只有一个方法，epoll操作过程有3个方法，分别是epoll_create()， epoll_ctl()，epoll_wait()。

### epoll_create
int epoll_create(int size)；
功能：用于创建一个epoll的句柄，size是指监听的描述符个数， 现在内核支持动态扩展，该值的意义仅仅是初次分配的fd个数，后面空间不够时会动态扩容。 当创建完epoll句柄后，占用一个fd值.

ls /proc/<pid>/fd/  //可通过终端执行，看到该fd
使用完epoll后，必须调用close()关闭，否则可能导致fd被耗尽。

### epoll_ctl
int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event)；
功能：用于对需要监听的文件描述符(fd)执行op操作，比如将fd加入到epoll句柄。

epfd：是epoll_create()的返回值；
op：表示op操作，用三个宏来表示，分别代表添加、删除和修改对fd的监听事件；
EPOLL_CTL_ADD(添加)
EPOLL_CTL_DEL(删除)
EPOLL_CTL_MOD（修改）
fd：需要监听的文件描述符；
epoll_event：需要监听的事件，struct epoll_event结构如下：
```
  struct epoll_event {
    __uint32_t events;  /* Epoll事件 */
    epoll_data_t data;  /*用户可用数据*/
  };
```
events可取值：(表示对应的文件描述符的操作)

EPOLLIN ：可读（包括对端SOCKET正常关闭）；
EPOLLOUT：可写；
EPOLLERR：错误；
EPOLLHUP：中断；
EPOLLPRI：高优先级的可读（这里应该表示有带外数据到来）；
EPOLLET： 将EPOLL设为边缘触发模式，这是相对于水平触发来说的。
EPOLLONESHOT：只监听一次事件，当监听完这次事件之后就不再监听该事件
### epoll_wait
int epoll_wait(int epfd, struct epoll_event * events, int maxevents, int timeout);
功能：等待事件的上报

epfd：等待epfd上的io事件，最多返回maxevents个事件；
events：用来从内核得到事件的集合；
maxevents：events数量，该maxevents值不能大于创建epoll_create()时的size；
timeout：超时时间（毫秒，0会立即返回）。
该函数返回需要处理的事件数目，如返回0表示已超时。


http://gityuan.com/2019/01/05/linux-poll-select/
