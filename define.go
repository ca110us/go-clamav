package goclamav

/*
#include <clamav.h>
#include <stdlib.h>
*/
import "C"
import "errors"

// ErrorCode models ClamAV errors
type ErrorCode C.cl_error_t

// scan options
type SCAN_OPTIONS struct {
	General   uint
	Parse     uint
	Heuristic uint
	Mail      uint
	Dev       uint
}

const (
	/* general */
	CL_SCAN_GENERAL_ALLMATCHES           = 0x1  /* scan in all-match mode */
	CL_SCAN_GENERAL_COLLECT_METADATA     = 0x2  /* collect metadata (--gen-json) */
	CL_SCAN_GENERAL_HEURISTICS           = 0x4  /* option to enable heuristic alerts */
	CL_SCAN_GENERAL_HEURISTIC_PRECEDENCE = 0x8  /* allow heuristic match to take precedence. */
	CL_SCAN_GENERAL_UNPRIVILEGED         = 0x10 /* scanner will not have read access to files. */

	/* parsing capabilities options */
	CL_SCAN_PARSE_ARCHIVE = 0x1
	CL_SCAN_PARSE_ELF     = 0x2
	CL_SCAN_PARSE_PDF     = 0x4
	CL_SCAN_PARSE_SWF     = 0x8
	CL_SCAN_PARSE_HWP3    = 0x10
	CL_SCAN_PARSE_XMLDOCS = 0x20
	CL_SCAN_PARSE_MAIL    = 0x40
	CL_SCAN_PARSE_OLE2    = 0x80
	CL_SCAN_PARSE_HTML    = 0x100
	CL_SCAN_PARSE_PE      = 0x200

	/* heuristic alerting options */
	CL_SCAN_HEURISTIC_BROKEN                  = 0x2    /* alert on broken PE and broken ELF files */
	CL_SCAN_HEURISTIC_EXCEEDS_MAX             = 0x4    /* alert when files exceed scan limits (filesize, max scansize, or max recursion depth) */
	CL_SCAN_HEURISTIC_PHISHING_SSL_MISMATCH   = 0x8    /* alert on SSL mismatches */
	CL_SCAN_HEURISTIC_PHISHING_CLOAK          = 0x10   /* alert on cloaked URLs in emails */
	CL_SCAN_HEURISTIC_MACROS                  = 0x20   /* alert on OLE2 files containing macros */
	CL_SCAN_HEURISTIC_ENCRYPTED_ARCHIVE       = 0x40   /* alert if archive is encrypted (rar, zip, etc) */
	CL_SCAN_HEURISTIC_ENCRYPTED_DOC           = 0x80   /* alert if a document is encrypted (pdf, docx, etc) */
	CL_SCAN_HEURISTIC_PARTITION_INTXN         = 0x100  /* alert if partition table size doesn't make sense */
	CL_SCAN_HEURISTIC_STRUCTURED              = 0x200  /* data loss prevention options, i.e. alert when detecting personal information */
	CL_SCAN_HEURISTIC_STRUCTURED_SSN_NORMAL   = 0x400  /* alert when detecting social security numbers */
	CL_SCAN_HEURISTIC_STRUCTURED_SSN_STRIPPED = 0x800  /* alert when detecting stripped social security numbers */
	CL_SCAN_HEURISTIC_STRUCTURED_CC           = 0x1000 /* alert when detecting credit card numbers */
	CL_SCAN_HEURISTIC_BROKEN_MEDIA            = 0x2000 /* alert if a file does not match the identified file format, works with JPEG, TIFF, GIF, PNG */

	/* mail scanning options */
	CL_SCAN_MAIL_PARTIAL_MESSAGE = 0x1

	/* dev options */
	CL_SCAN_DEV_COLLECT_SHA              = 0x1 /* Enables hash output in sha-collect builds - for internal use only */
	CL_SCAN_DEV_COLLECT_PERFORMANCE_INFO = 0x2 /* collect performance timings */

	/* cl_countsigs options */
	CL_COUNTSIGS_OFFICIAL   = 0x1
	CL_COUNTSIGS_UNOFFICIAL = 0x2
	CL_COUNTSIGS_ALL        = (CL_COUNTSIGS_OFFICIAL | CL_COUNTSIGS_UNOFFICIAL)
)

// Fmap models in-memory files
type Fmap C.cl_fmap_t

const (
	/* libclamav specific */
	CL_CLEAN    ErrorCode = C.CL_CLEAN
	CL_SUCCESS  ErrorCode = C.CL_SUCCESS
	CL_VIRUS    ErrorCode = C.CL_VIRUS
	CL_ENULLARG ErrorCode = C.CL_ENULLARG
	CL_EARG     ErrorCode = C.CL_EARG
	CL_EMALFDB  ErrorCode = C.CL_EMALFDB
	CL_ECVD     ErrorCode = C.CL_ECVD
	CL_EVERIFY  ErrorCode = C.CL_EVERIFY
	CL_EUNPACK  ErrorCode = C.CL_EUNPACK

	/* I/O and memory errors */
	CL_EOPEN    ErrorCode = C.CL_EOPEN
	CL_ECREAT   ErrorCode = C.CL_ECREAT
	CL_EUNLINK  ErrorCode = C.CL_EUNLINK
	CL_ESTAT    ErrorCode = C.CL_ESTAT
	CL_EREAD    ErrorCode = C.CL_EREAD
	CL_ESEEK    ErrorCode = C.CL_ESEEK
	CL_EWRITE   ErrorCode = C.CL_EWRITE
	CL_EDUP     ErrorCode = C.CL_EDUP
	CL_EACCES   ErrorCode = C.CL_EACCES
	CL_ETMPFILE ErrorCode = C.CL_ETMPFILE
	CL_ETMPDIR  ErrorCode = C.CL_ETMPDIR
	CL_EMAP     ErrorCode = C.CL_EMAP
	CL_EMEM     ErrorCode = C.CL_EMEM
	CL_ETIMEOUT ErrorCode = C.CL_ETIMEOUT

	/* internal (not reported outside libclamav) */
	CL_BREAK              ErrorCode = C.CL_BREAK
	CL_EMAXREC            ErrorCode = C.CL_EMAXREC
	CL_EMAXSIZE           ErrorCode = C.CL_EMAXSIZE
	CL_EMAXFILES          ErrorCode = C.CL_EMAXFILES
	CL_EFORMAT            ErrorCode = C.CL_EFORMAT
	CL_EPARSE             ErrorCode = C.CL_EPARSE
	CL_EBYTECODE          ErrorCode = C.CL_EBYTECODE          /* may be reported in testmode */
	CL_EBYTECODE_TESTFAIL ErrorCode = C.CL_EBYTECODE_TESTFAIL /* may be reported in testmode */

	/* c4w error codes */
	CL_ELOCK  ErrorCode = C.CL_ELOCK
	CL_EBUSY  ErrorCode = C.CL_EBUSY
	CL_ESTATE ErrorCode = C.CL_ESTATE

	CL_VERIFIED ErrorCode = C.CL_VERIFIED /* The binary has been deemed trusted */
	CL_ERROR    ErrorCode = C.CL_ERROR    /* Unspecified / generic error */

	/* no error codes below this line please */
	CL_ELAST_ERROR ErrorCode = C.CL_ELAST_ERROR
)

// CL_INIT_DEFAULT is a macro that can be passed to cl_init() representing the default initialization settings
const CL_INIT_DEFAULT C.uint = C.CL_INIT_DEFAULT

// Wraps the corresponding error message
func Strerr(code ErrorCode) error {
	err := errors.New(C.GoString(C.cl_strerror(C.cl_error_t(code))))
	return err
}

/* db options */
// clang-format off
type DBOptions uint

const (
	CL_DB_PHISHING          DBOptions = C.CL_DB_PHISHING
	CL_DB_PHISHING_URLS     DBOptions = C.CL_DB_PHISHING_URLS
	CL_DB_PUA               DBOptions = C.CL_DB_PUA
	CL_DB_CVDNOTMP          DBOptions = C.CL_DB_CVDNOTMP /* obsolete */
	CL_DB_OFFICIAL          DBOptions = C.CL_DB_OFFICIAL /* internal */
	CL_DB_PUA_MODE          DBOptions = C.CL_DB_PUA_MODE
	CL_DB_PUA_INCLUDE       DBOptions = C.CL_DB_PUA_INCLUDE
	CL_DB_PUA_EXCLUDE       DBOptions = C.CL_DB_PUA_EXCLUDE
	CL_DB_COMPILED          DBOptions = C.CL_DB_COMPILED  /* internal */
	CL_DB_DIRECTORY         DBOptions = C.CL_DB_DIRECTORY /* internal */
	CL_DB_OFFICIAL_ONLY     DBOptions = C.CL_DB_OFFICIAL_ONLY
	CL_DB_BYTECODE          DBOptions = C.CL_DB_BYTECODE
	CL_DB_SIGNED            DBOptions = C.CL_DB_SIGNED            /* internal */
	CL_DB_BYTECODE_UNSIGNED DBOptions = C.CL_DB_BYTECODE_UNSIGNED /* Caution: You should never run bytecode signatures from untrusted sources. Doing so may result in arbitrary code execution. */
	CL_DB_UNSIGNED          DBOptions = C.CL_DB_UNSIGNED          /* internal */
	CL_DB_BYTECODE_STATS    DBOptions = C.CL_DB_BYTECODE_STATS
	CL_DB_ENHANCED          DBOptions = C.CL_DB_ENHANCED
	CL_DB_PCRE_STATS        DBOptions = C.CL_DB_PCRE_STATS
	CL_DB_YARA_EXCLUDE      DBOptions = C.CL_DB_YARA_EXCLUDE
	CL_DB_YARA_ONLY         DBOptions = C.CL_DB_YARA_ONLY
)

// EngineField selects a particular engine settings field
type EngineField C.enum_cl_engine_field

// Engine settings
const (
	CL_ENGINE_MAX_SCANSIZE        EngineField = C.CL_ENGINE_MAX_SCANSIZE        /* uint64_t */
	CL_ENGINE_MAX_FILESIZE        EngineField = C.CL_ENGINE_MAX_FILESIZE        /* uint64_t */
	CL_ENGINE_MAX_RECURSION       EngineField = C.CL_ENGINE_MAX_RECURSION       /* uint32_t */
	CL_ENGINE_MAX_FILES           EngineField = C.CL_ENGINE_MAX_FILES           /* uint32_t */
	CL_ENGINE_MIN_CC_COUNT        EngineField = C.CL_ENGINE_MIN_CC_COUNT        /* uint32_t */
	CL_ENGINE_MIN_SSN_COUNT       EngineField = C.CL_ENGINE_MIN_SSN_COUNT       /* uint32_t */
	CL_ENGINE_PUA_CATEGORIES      EngineField = C.CL_ENGINE_PUA_CATEGORIES      /* (char *) */
	CL_ENGINE_DB_OPTIONS          EngineField = C.CL_ENGINE_DB_OPTIONS          /* uint32_t */
	CL_ENGINE_DB_VERSION          EngineField = C.CL_ENGINE_DB_VERSION          /* uint32_t */
	CL_ENGINE_DB_TIME             EngineField = C.CL_ENGINE_DB_TIME             /* time_t */
	CL_ENGINE_AC_ONLY             EngineField = C.CL_ENGINE_AC_ONLY             /* uint32_t */
	CL_ENGINE_AC_MINDEPTH         EngineField = C.CL_ENGINE_AC_MINDEPTH         /* uint32_t */
	CL_ENGINE_AC_MAXDEPTH         EngineField = C.CL_ENGINE_AC_MAXDEPTH         /* uint32_t */
	CL_ENGINE_TMPDIR              EngineField = C.CL_ENGINE_TMPDIR              /* (char *) */
	CL_ENGINE_KEEPTMP             EngineField = C.CL_ENGINE_KEEPTMP             /* uint32_t */
	CL_ENGINE_BYTECODE_SECURITY   EngineField = C.CL_ENGINE_BYTECODE_SECURITY   /* uint32_t */
	CL_ENGINE_BYTECODE_TIMEOUT    EngineField = C.CL_ENGINE_BYTECODE_TIMEOUT    /* uint32_t */
	CL_ENGINE_BYTECODE_MODE       EngineField = C.CL_ENGINE_BYTECODE_MODE       /* uint32_t */
	CL_ENGINE_MAX_EMBEDDEDPE      EngineField = C.CL_ENGINE_MAX_EMBEDDEDPE      /* uint64_t */
	CL_ENGINE_MAX_HTMLNORMALIZE   EngineField = C.CL_ENGINE_MAX_HTMLNORMALIZE   /* uint64_t */
	CL_ENGINE_MAX_HTMLNOTAGS      EngineField = C.CL_ENGINE_MAX_HTMLNOTAGS      /* uint64_t */
	CL_ENGINE_MAX_SCRIPTNORMALIZE EngineField = C.CL_ENGINE_MAX_SCRIPTNORMALIZE /* uint64_t */
	CL_ENGINE_MAX_ZIPTYPERCG      EngineField = C.CL_ENGINE_MAX_ZIPTYPERCG      /* uint64_t */
	CL_ENGINE_FORCETODISK         EngineField = C.CL_ENGINE_FORCETODISK         /* uint32_t */
	CL_ENGINE_DISABLE_CACHE       EngineField = C.CL_ENGINE_DISABLE_CACHE       /* uint32_t */
	CL_ENGINE_DISABLE_PE_STATS    EngineField = C.CL_ENGINE_DISABLE_PE_STATS    /* uint32_t */
	CL_ENGINE_STATS_TIMEOUT       EngineField = C.CL_ENGINE_STATS_TIMEOUT       /* uint32_t */
	CL_ENGINE_MAX_PARTITIONS      EngineField = C.CL_ENGINE_MAX_PARTITIONS      /* uint32_t */
	CL_ENGINE_MAX_ICONSPE         EngineField = C.CL_ENGINE_MAX_ICONSPE         /* uint32_t */
	CL_ENGINE_MAX_RECHWP3         EngineField = C.CL_ENGINE_MAX_RECHWP3         /* uint32_t */
	CL_ENGINE_MAX_SCANTIME        EngineField = C.CL_ENGINE_MAX_SCANTIME        /* uint32_t */
	CL_ENGINE_PCRE_MATCH_LIMIT    EngineField = C.CL_ENGINE_PCRE_MATCH_LIMIT    /* uint64_t */
	CL_ENGINE_PCRE_RECMATCH_LIMIT EngineField = C.CL_ENGINE_PCRE_RECMATCH_LIMIT /* uint64_t */
	CL_ENGINE_PCRE_MAX_FILESIZE   EngineField = C.CL_ENGINE_PCRE_MAX_FILESIZE   /* uint64_t */
	CL_ENGINE_DISABLE_PE_CERTS    EngineField = C.CL_ENGINE_DISABLE_PE_CERTS    /* uint32_t */
	CL_ENGINE_PE_DUMPCERTS        EngineField = C.CL_ENGINE_PE_DUMPCERTS        /* uint32_t */
)
