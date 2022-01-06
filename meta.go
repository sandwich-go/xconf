package xconf

const (
	MetaKeyFiles                            = "xconf_files"
	MetaKeyGrayLabel                        = "xconf_gray_rule_label"
	MetaKeyInheritFiles                     = "xconf_inherit_files"
	MetaKeyLatestHash                       = "xconf_latest_hash"
	MetaKeyInheritFilesDeprecatedFromGoconf = "inherit_files"
)

const (
	// HashPrefix hash字段前缀
	HashPrefix = "xconf@"
	// DefaultInvalidHashString 默认hash值
	DefaultInvalidHashString = HashPrefix + "hash_invalid"
)

var MetaKeyList = []string{
	MetaKeyInheritFilesDeprecatedFromGoconf,
	MetaKeyInheritFiles,
	MetaKeyLatestHash,
	MetaKeyGrayLabel,
	MetaKeyFiles}
