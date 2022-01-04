package xconf

const MetaKeyFiles = "xconf_files"
const MetaKeyInheritFiles = "xconf_inherit_files"
const MetaKeyLatestHash = "xconf_latest_hash"
const HashPrefix = "xconf@"
const DefaultInvalidHashString = HashPrefix + "hash_invalid"

const MetaKeyInheritFilesDeprecatedFromGoconf = "inherit_files"

var MetaKeyList = []string{
	MetaKeyInheritFilesDeprecatedFromGoconf,
	MetaKeyInheritFiles,
	MetaKeyLatestHash,
	MetaKeyFiles}
