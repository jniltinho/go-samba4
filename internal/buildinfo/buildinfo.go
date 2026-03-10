package buildinfo

// Version and BuildDate are injected at build time via ldflags:
//
//	-X 'go-samba4/internal/buildinfo.Version=v1.x.x'
//	-X 'go-samba4/internal/buildinfo.BuildDate=2006-01-02'
var (
	Version   = "v1.0.0"
	BuildDate = "Unknown"
)
