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

	// MetaKeyIgnoreFields 元数据，配置文件内使用，用于声明字段路径白名单
	// 命中的字段即使不在目标 struct 中声明，xconf 解析时也不会报"未识别字段"错误
	// 适用于插件/业务自定义字段、平滑迁移等场景
	// 支持精确匹配或前缀匹配（按 . 边界），如 "plugin.foo" 同时命中 "plugin.foo"、"plugin.foo.bar"，但不命中 "plugin.foobar"
	// 在 yaml 中示例:
	//   xconf_ignore_fields:
	//     - csharp.gen_csharp_proto_db
	//     - my_plugin
	// 多文件继承时，多个文件的白名单会做并集去重合并
	MetaKeyIgnoreFields = "xconf_ignore_fields"
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
	MetaKeyFlagFiles,
	MetaKeyIgnoreFields}
