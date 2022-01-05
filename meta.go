package xconf

const (
	MetaKeyFiles                            = "xconf_files"
	MetaKeyGrayLabel                        = "xconf_gray_rule_label"
	MetaKeyInheritFiles                     = "xconf_inherit_files"
	MetaKeyLatestHash                       = "xconf_latest_hash"
	MetaKeyInheritFilesDeprecatedFromGoconf = "inherit_files"
)

const (
	HashPrefix               = "xconf@"
	DefaultInvalidHashString = HashPrefix + "hash_invalid"
)

var MetaKeyList = []string{
	MetaKeyInheritFilesDeprecatedFromGoconf,
	MetaKeyInheritFiles,
	MetaKeyLatestHash,
	MetaKeyGrayLabel,
	MetaKeyFiles}
