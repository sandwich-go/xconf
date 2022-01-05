package xconf

const (
	MetaKeyFiles             = "xconf_files"
	MetaKeyGrayLabel         = "xconf_gray_rule_label"
	MetaKeyInheritFiles      = "xconf_inherit_files"
	MetaKeyLatestHash        = "xconf_latest_hash"
	HashPrefix               = "xconf@"
	DefaultInvalidHashString = HashPrefix + "hash_invalid"
)

const MetaKeyInheritFilesDeprecatedFromGoconf = "inherit_files"

var MetaKeyList = []string{
	MetaKeyInheritFilesDeprecatedFromGoconf,
	MetaKeyInheritFiles,
	MetaKeyLatestHash,
	MetaKeyGrayLabel,
	MetaKeyFiles}
