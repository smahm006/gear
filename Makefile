VERSION := "0.0.1"
BUILD := `git rev-parse --short HEAD`
VERSIONFILE := src/cmd/version.go

gensrc:
	rm -f $(VERSIONFILE)
	@echo "package cmd" > $(VERSIONFILE)
	@echo "const (" >> $(VERSIONFILE)
	@echo "  VERSION = \"$(VERSION)\"" >> $(VERSIONFILE)
	@echo "  BUILD = \"$(BUILD)\"" >> $(VERSIONFILE)
	@echo ")" >> $(VERSIONFILE)
