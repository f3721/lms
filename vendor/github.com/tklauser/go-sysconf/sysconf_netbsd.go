// Copyright 2018 Tobias Klauser. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sysconf

import (
	"sync"

	"golang.org/x/sys/unix"
)

const (
	_HOST_NAME_MAX  = _MAXHOSTNAMELEN
	_LOGIN_NAME_MAX = _MAXLOGNAME + 1
	_SYMLOOP_MAX    = _MAXSYMLINKS

	_POSIX2_C_DEV = -1
	_POSIX2_UPE   = -1
)

var (
	clktck     int64
	clktckOnce sync.Once
)

func sysconfPOSIX(name int) (int64, error) {
	// NetBSD does not define all _POSIX_* values used in sysconf_posix.go
	// Handle the supported ones here.
	switch name {
	case SC_SHELL:
		return _POSIX_SHELL, nil
	case SC_VERSION:
		return _POSIX_VERSION, nil
	}

	return -1, errInvalid
}

func sysconf(name int) (int64, error) {
	// NetBSD uses sysctl to get some of these values. For the user.* namespace,
	// calls get handled by user_sysctl in /usr/src/lib/libc/gen/sysctl.c
	// Duplicate the relevant values here.

	switch name {
	case SC_ARG_MAX:
		return sysctl32("kern.argmax"), nil
	case SC_CHILD_MAX:
		var rlim unix.Rlimit
		if err := unix.Getrlimit(unix.RLIMIT_NPROC, &rlim); err == nil {
			if rlim.Cur != unix.RLIM_INFINITY {
				return int64(rlim.Cur), nil
			}
		}
		return -1, nil
	case SC_STREAM_MAX:
		// sysctl("user.stream_max")
		return _FOPEN_MAX, nil
	case SC_TTY_NAME_MAX:
		return pathconf(_PATH_DEV, _PC_NAME_MAX), nil
	case SC_CLK_TCK:
		clktckOnce.Do(func() {
			clktck = -1
			if ci, err := unix.SysctlClockinfo("kern.clockrate"); err == nil {
				clktck = int64(ci.Hz)
			}
		})
		return clktck, nil
	case SC_NGROUPS_MAX:
		return sysctl32("kern.ngroups"), nil
	case SC_JOB_CONTROL:
		return sysctl32("kern.job_control"), nil
	case SC_OPEN_MAX:
		var rlim unix.Rlimit
		if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &rlim); err == nil {
			return int64(rlim.Cur), nil
		}
		return -1, nil
	case SC_TZNAME_MAX:
		// sysctl("user.tzname_max")
		return _NAME_MAX, nil

	// 1003.1b
	case SC_FSYNC:
		return sysctl32("kern.fsync"), nil
	case SC_MAPPED_FILES:
		return sysctl32("kern.mapped_files"), nil
	case SC_MONOTONIC_CLOCK:
		return sysctl32("kern.monotonic_clock"), nil
	case SC_SEMAPHORES:
		return sysctl32("kern.posix_semaphores"), nil
	case SC_TIMERS:
		return sysctl32("kern.posix_timers"), nil

	// 1003.1c
	case SC_LOGIN_NAME_MAX:
		return sysctl32("kern.login_name_max"), nil
	case SC_THREADS:
		return sysctl32("kern.posix_threads"), nil

	// 1003.1j
	case SC_BARRIERS:
		return sysctl32("kern.posix_barriers"), nil

	// 1003.2
	case SC_2_VERSION:
		// sysctl("user.posix2_version")
		return _POSIX2_VERSION, nil
	case SC_2_UPE:
		// sysctl("user.posix2_upe")
		return _POSIX2_UPE, nil

	// XPG 4.2
	case SC_IOV_MAX:
		return sysctl32("kern.iov_max"), nil

	// 1003.1-2001, XSI Option Group
	case SC_AIO_LISTIO_MAX:
		return sysctl32("kern.aio_listio_max"), nil
	case SC_AIO_MAX:
		return sysctl32("kern.aio_max"), nil
	case SC_ASYNCHRONOUS_IO:
		return sysctl32("kern.posix_aio"), nil
	case SC_MQ_OPEN_MAX:
		return sysctl32("kern.mqueue.mq_open_max"), nil
	case SC_MQ_PRIO_MAX:
		return sysctl32("kern.mqueue.mq_prio_max"), nil
	case SC_ATEXIT_MAX:
		// sysctl("user.atexit_max")
		return -1, nil // TODO

	// Extensions
	case SC_NPROCESSORS_CONF:
		return sysctl32("hw.ncpu"), nil
	case SC_NPROCESSORS_ONLN:
		return sysctl32("hw.ncpuonline"), nil

	// Linux/Solaris
	case SC_PHYS_PAGES:
		return sysctl64("hw.physmem64") / int64(unix.Getpagesize()), nil

	// Native
	case SC_THREAD_DESTRUCTOR_ITERATIONS:
		return _POSIX_THREAD_DESTRUCTOR_ITERATIONS, nil
	case SC_THREAD_KEYS_MAX:
		return _POSIX_THREAD_KEYS_MAX, nil
	case SC_THREAD_STACK_MIN:
		return int64(unix.Getpagesize()), nil
	case SC_THREAD_THREADS_MAX:
		return sysctl32("kern.maxproc"), nil
	}

	return sysconfGeneric(name)
}
