# aptlib [![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-red.svg)](https://choosealicense.com/licenses/apache-2.0)

A library for working with apt implemented in Go.

## To Do

- [x] Correctly parse sources.list files.
- [ ] Support all valid source options.
	- [ ] arch
		- [x] single
		- [ ] multiple
	- [ ] lang
		- [ ] single
		- [ ] multiple
	- [ ] target
	- [ ] pdiffs
	- [ ] by-hash
	- [ ] allow-insecure
	- [ ] allow-weak
	- [ ] allow-downgrade-to-insecure
	- [ ] trusted
	- [ ] signed-by
	- [ ] check-valid-until
	- [ ] valid-until-min
	- [ ] check-date
	- [ ] date-max-future
	- [ ] inrelease-path
- [ ] Support all valid source types.
	- [x] HTTP(S)
	- [ ] FTP
	- [ ] SSH
	- [ ] Copy
- [x] Download package indices.
- [x] UnGzip indices.
- [ ] PGP signature verification.
- [ ] Parse package indices for package information.
- [ ] Download individual packages.
- [ ] Package index searching.
- [ ] Enumerate upgradable packages.
- [ ] Support local repository mirroring.
