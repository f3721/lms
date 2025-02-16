// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs sysconf_defs_openbsd.go

//go:build openbsd
// +build openbsd

package sysconf

const (
	SC_AIO_LISTIO_MAX               = 0x2a
	SC_AIO_MAX                      = 0x2b
	SC_AIO_PRIO_DELTA_MAX           = 0x2c
	SC_ARG_MAX                      = 0x1
	SC_ATEXIT_MAX                   = 0x2e
	SC_BC_BASE_MAX                  = 0x9
	SC_BC_DIM_MAX                   = 0xa
	SC_BC_SCALE_MAX                 = 0xb
	SC_BC_STRING_MAX                = 0xc
	SC_CHILD_MAX                    = 0x2
	SC_CLK_TCK                      = 0x3
	SC_COLL_WEIGHTS_MAX             = 0xd
	SC_DELAYTIMER_MAX               = 0x32
	SC_EXPR_NEST_MAX                = 0xe
	SC_GETGR_R_SIZE_MAX             = 0x64
	SC_GETPW_R_SIZE_MAX             = 0x65
	SC_HOST_NAME_MAX                = 0x21
	SC_IOV_MAX                      = 0x33
	SC_LINE_MAX                     = 0xf
	SC_LOGIN_NAME_MAX               = 0x66
	SC_MQ_OPEN_MAX                  = 0x3a
	SC_MQ_PRIO_MAX                  = 0x3b
	SC_NGROUPS_MAX                  = 0x4
	SC_OPEN_MAX                     = 0x5
	SC_PAGE_SIZE                    = 0x1c
	SC_PAGESIZE                     = 0x1c
	SC_THREAD_DESTRUCTOR_ITERATIONS = 0x50
	SC_THREAD_KEYS_MAX              = 0x51
	SC_THREAD_STACK_MIN             = 0x59
	SC_THREAD_THREADS_MAX           = 0x5a
	SC_RE_DUP_MAX                   = 0x10
	SC_SEM_NSEMS_MAX                = 0x1f
	SC_SEM_VALUE_MAX                = 0x20
	SC_SIGQUEUE_MAX                 = 0x46
	SC_STREAM_MAX                   = 0x1a
	SC_SYMLOOP_MAX                  = 0x4c
	SC_TIMER_MAX                    = 0x5d
	SC_TTY_NAME_MAX                 = 0x6b
	SC_TZNAME_MAX                   = 0x1b

	SC_ADVISORY_INFO              = 0x29
	SC_ASYNCHRONOUS_IO            = 0x2d
	SC_BARRIERS                   = 0x2f
	SC_CLOCK_SELECTION            = 0x30
	SC_CPUTIME                    = 0x31
	SC_FSYNC                      = 0x1d
	SC_IPV6                       = 0x34
	SC_JOB_CONTROL                = 0x6
	SC_MAPPED_FILES               = 0x35
	SC_MEMLOCK                    = 0x36
	SC_MEMLOCK_RANGE              = 0x37
	SC_MEMORY_PROTECTION          = 0x38
	SC_MESSAGE_PASSING            = 0x39
	SC_MONOTONIC_CLOCK            = 0x22
	SC_PRIORITIZED_IO             = 0x3c
	SC_PRIORITY_SCHEDULING        = 0x3d
	SC_RAW_SOCKETS                = 0x3e
	SC_READER_WRITER_LOCKS        = 0x3f
	SC_REALTIME_SIGNALS           = 0x40
	SC_REGEXP                     = 0x41
	SC_SAVED_IDS                  = 0x7
	SC_SEMAPHORES                 = 0x43
	SC_SHARED_MEMORY_OBJECTS      = 0x44
	SC_SHELL                      = 0x45
	SC_SPAWN                      = 0x47
	SC_SPIN_LOCKS                 = 0x48
	SC_SPORADIC_SERVER            = 0x49
	SC_SS_REPL_MAX                = 0x4a
	SC_SYNCHRONIZED_IO            = 0x4b
	SC_THREAD_ATTR_STACKADDR      = 0x4d
	SC_THREAD_ATTR_STACKSIZE      = 0x4e
	SC_THREAD_CPUTIME             = 0x4f
	SC_THREAD_PRIO_INHERIT        = 0x52
	SC_THREAD_PRIO_PROTECT        = 0x53
	SC_THREAD_PRIORITY_SCHEDULING = 0x54
	SC_THREAD_PROCESS_SHARED      = 0x55
	SC_THREAD_ROBUST_PRIO_INHERIT = 0x56
	SC_THREAD_ROBUST_PRIO_PROTECT = 0x57
	SC_THREAD_SAFE_FUNCTIONS      = 0x67
	SC_THREAD_SPORADIC_SERVER     = 0x58
	SC_THREADS                    = 0x5b
	SC_TIMEOUTS                   = 0x5c
	SC_TIMERS                     = 0x5e
	SC_TRACE                      = 0x5f
	SC_TRACE_EVENT_FILTER         = 0x60
	SC_TRACE_EVENT_NAME_MAX       = 0x61
	SC_TRACE_INHERIT              = 0x62
	SC_TRACE_LOG                  = 0x63
	SC_TRACE_NAME_MAX             = 0x68
	SC_TRACE_SYS_MAX              = 0x69
	SC_TRACE_USER_EVENT_MAX       = 0x6a
	SC_TYPED_MEMORY_OBJECTS       = 0x6c
	SC_VERSION                    = 0x8

	SC_V7_ILP32_OFF32  = 0x71
	SC_V7_ILP32_OFFBIG = 0x72
	SC_V7_LP64_OFF64   = 0x73
	SC_V7_LPBIG_OFFBIG = 0x74

	SC_V6_ILP32_OFF32  = 0x6d
	SC_V6_ILP32_OFFBIG = 0x6e
	SC_V6_LP64_OFF64   = 0x6f
	SC_V6_LPBIG_OFFBIG = 0x70

	SC_2_C_BIND         = 0x12
	SC_2_C_DEV          = 0x13
	SC_2_CHAR_TERM      = 0x14
	SC_2_FORT_DEV       = 0x15
	SC_2_FORT_RUN       = 0x16
	SC_2_LOCALEDEF      = 0x17
	SC_2_PBS            = 0x23
	SC_2_PBS_ACCOUNTING = 0x24
	SC_2_PBS_CHECKPOINT = 0x25
	SC_2_PBS_LOCATE     = 0x26
	SC_2_PBS_MESSAGE    = 0x27
	SC_2_PBS_TRACK      = 0x28
	SC_2_SW_DEV         = 0x18
	SC_2_UPE            = 0x19
	SC_2_VERSION        = 0x11

	SC_XOPEN_CRYPT            = 0x75
	SC_XOPEN_ENH_I18N         = 0x76
	SC_XOPEN_REALTIME         = 0x78
	SC_XOPEN_REALTIME_THREADS = 0x79
	SC_XOPEN_SHM              = 0x1e
	SC_XOPEN_STREAMS          = 0x7a
	SC_XOPEN_UNIX             = 0x7b
	SC_XOPEN_UUCP             = 0x7c
	SC_XOPEN_VERSION          = 0x7d

	SC_AVPHYS_PAGES     = 0x1f5
	SC_PHYS_PAGES       = 0x1f4
	SC_NPROCESSORS_CONF = 0x1f6
	SC_NPROCESSORS_ONLN = 0x1f7
)

const (
	_HOST_NAME_MAX                 = 0xff
	_IOV_MAX                       = 0x400
	_LOGIN_NAME_MAX                = 0x20
	_PTHREAD_DESTRUCTOR_ITERATIONS = 0x4
	_PTHREAD_KEYS_MAX              = 0x100
	_PTHREAD_STACK_MIN             = 0x1000
	_PTHREAD_THREADS_MAX           = 0xffffffffffffffff
	_SEM_VALUE_MAX                 = 0xffffffff
	_SYMLOOP_MAX                   = 0x20
	_TTY_NAME_MAX                  = 0x104

	_GR_BUF_LEN = 0xa40
	_PW_BUF_LEN = 0x400

	_CLK_TCK = 0x64

	_POSIX_ADVISORY_INFO              = -0x1
	_POSIX_ARG_MAX                    = 0x1000
	_POSIX_ASYNCHRONOUS_IO            = -0x1
	_POSIX_BARRIERS                   = 0x30db0
	_POSIX_CHILD_MAX                  = 0x19
	_POSIX_CLOCK_SELECTION            = -0x1
	_POSIX_CPUTIME                    = 0x31069
	_POSIX_FSYNC                      = 0x30db0
	_POSIX_IPV6                       = 0x0
	_POSIX_JOB_CONTROL                = 0x1
	_POSIX_MAPPED_FILES               = 0x30db0
	_POSIX_MEMLOCK                    = 0x30db0
	_POSIX_MEMLOCK_RANGE              = 0x30db0
	_POSIX_MEMORY_PROTECTION          = 0x30db0
	_POSIX_MESSAGE_PASSING            = -0x1
	_POSIX_MONOTONIC_CLOCK            = 0x30db0
	_POSIX_PRIORITIZED_IO             = -0x1
	_POSIX_PRIORITY_SCHEDULING        = -0x1
	_POSIX_RAW_SOCKETS                = 0x30db0
	_POSIX_READER_WRITER_LOCKS        = 0x30db0
	_POSIX_REALTIME_SIGNALS           = -0x1
	_POSIX_REGEXP                     = 0x1
	_POSIX_SAVED_IDS                  = 0x1
	_POSIX_SEMAPHORES                 = 0x30db0
	_POSIX_SHARED_MEMORY_OBJECTS      = 0x31069
	_POSIX_SHELL                      = 0x1
	_POSIX_SPAWN                      = 0x30db0
	_POSIX_SPIN_LOCKS                 = 0x30db0
	_POSIX_SPORADIC_SERVER            = -0x1
	_POSIX_SYNCHRONIZED_IO            = -0x1
	_POSIX_THREAD_ATTR_STACKADDR      = 0x30db0
	_POSIX_THREAD_ATTR_STACKSIZE      = 0x30db0
	_POSIX_THREAD_CPUTIME             = 0x31069
	_POSIX_THREAD_KEYS_MAX            = 0x80
	_POSIX_THREAD_PRIO_INHERIT        = -0x1
	_POSIX_THREAD_PRIO_PROTECT        = -0x1
	_POSIX_THREAD_PRIORITY_SCHEDULING = -0x1
	_POSIX_THREAD_PROCESS_SHARED      = -0x1
	_POSIX_THREAD_ROBUST_PRIO_INHERIT = -0x1
	_POSIX_THREAD_ROBUST_PRIO_PROTECT = -0x1
	_POSIX_THREAD_SAFE_FUNCTIONS      = 0x30db0
	_POSIX_THREAD_SPORADIC_SERVER     = -0x1
	_POSIX_THREADS                    = 0x30db0
	_POSIX_TIMERS                     = -0x1
	_POSIX_TIMEOUTS                   = 0x30db0
	_POSIX_TRACE                      = -0x1
	_POSIX_TYPED_MEMORY_OBJECTS       = -0x1
	_POSIX_VERSION                    = 0x31069

	_POSIX_V7_ILP32_OFF32  = -0x1
	_POSIX_V7_ILP32_OFFBIG = 0x0
	_POSIX_V7_LP64_OFF64   = 0x0
	_POSIX_V7_LPBIG_OFFBIG = 0x0

	_POSIX_V6_ILP32_OFF32  = -0x1
	_POSIX_V6_ILP32_OFFBIG = 0x0
	_POSIX_V6_LP64_OFF64   = 0x0
	_POSIX_V6_LPBIG_OFFBIG = 0x0

	_POSIX2_C_BIND    = 0x30db0
	_POSIX2_C_DEV     = -0x1
	_POSIX2_CHAR_TERM = 0x1
	_POSIX2_LOCALEDEF = -0x1
	_POSIX2_PBS       = -0x1
	_POSIX2_SW_DEV    = 0x30db0
	_POSIX2_UPE       = 0x30db0
	_POSIX2_VERSION   = 0x31069

	_XOPEN_CRYPT            = 0x1
	_XOPEN_ENH_I18N         = -0x1
	_XOPEN_REALTIME         = -0x1
	_XOPEN_REALTIME_THREADS = -0x1
	_XOPEN_SHM              = 0x1
	_XOPEN_STREAMS          = -0x1
	_XOPEN_UNIX             = -0x1
	_XOPEN_UUCP             = -0x1

	_FOPEN_MAX  = 0x14
	_NAME_MAX   = 0xff
	_RE_DUP_MAX = 0xff

	_BC_BASE_MAX      = 0x7fffffff
	_BC_DIM_MAX       = 0xffff
	_BC_SCALE_MAX     = 0x7fffffff
	_BC_STRING_MAX    = 0x7fffffff
	_COLL_WEIGHTS_MAX = 0x2
	_EXPR_NEST_MAX    = 0x20
	_LINE_MAX         = 0x800

	_SHRT_MAX = 0x7fff

	_PATH_ZONEINFO = "/usr/share/zoneinfo"
)

const (
	_CHAR_BIT = 0x8

	_INT_MAX = 0x7fffffff

	sizeofOffT = 0x8
)
