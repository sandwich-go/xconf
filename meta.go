package xconf

const (
	// MetaKeyFlagFiles 元数据，flag中使用，用于通过flag指定需加载的配置文件列表
	// 多个文件以,分割，如server --xconf_flag_files=base.yaml,testflight.yaml
	MetaKeyFlagFiles = "xconf_flag_files"

	// MetaKeyInheritFiles 元数据，配置文件内使用，用于指定继承的文件
	// 如toml中配置：xconf_inherit_files=[etcd.yaml,production.yaml],则当前配置会继承etcd.yaml,production.yaml文件
	MetaKeyInheritFiles = "xconf_inherit_files"

	// MetaKeyGrayLabel 元数据，灰度发布支持，发布配置的的时候指定配置生次奥的label
	// xconf运行时可以通过WithAppLabelList指定label，当MetaKeyGrayLabel不为空且至少含有一个AppLabelList中的label时配置会应用到当前实例
	MetaKeyGrayLabel = "xconf_gray_rule_label"

	// MetaKeyLatestHash 元数据，预留用于配置版本的比对、校验
	MetaKeyLatestHash = "xconf_latest_hash"

	// MetaKeyInheritFilesDeprecatedFromGoconf 同MetaKeyInheritFiles，兼容goconf
	MetaKeyInheritFilesDeprecatedFromGoconf = "inherit_files"
)

const (
	// HashPrefix hash字段前缀
	HashPrefix = "xconf@"
	// DefaultInvalidHashString 默认hash值
	DefaultInvalidHashString = HashPrefix + "hash_invalid"
)

var metaKeyList = []string{
	MetaKeyInheritFilesDeprecatedFromGoconf,
	MetaKeyInheritFiles,
	MetaKeyLatestHash,
	MetaKeyGrayLabel,
	MetaKeyFlagFiles}
