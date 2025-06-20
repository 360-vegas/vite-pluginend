// 语言类型定义
export interface Language {
  name: string
  country: string
  flag: string
  prompt?: string
}

export interface LanguageGroup {
  [key: string]: Language
}

export interface LanguageGroups {
  [region: string]: LanguageGroup
}

// 语言配置数据
export const languageGroups: LanguageGroups = {
  '东亚': {
    'zh-CN': {
      name: '简体中文',
      country: '中国大陆',
      flag: 'cn',
      prompt: 'You must write the entire article in Simplified Chinese (简体中文). Use appropriate Chinese punctuation and expressions common in mainland China.'
    },
    'zh-TW': {
      name: '繁体中文(台湾)',
      country: '中国台湾',
      flag: 'tw',
      prompt: 'You must write the entire article in Traditional Chinese (繁體中文). Use appropriate Chinese punctuation and expressions common in Taiwan.'
    },
    'zh-HK': {
      name: '繁体中文(香港)',
      country: '中国香港',
      flag: 'hk',
      prompt: 'You must write the entire article in Traditional Chinese (繁體中文). Use appropriate Chinese punctuation and expressions common in Hong Kong.'
    },
    'zh-MO': {
      name: '繁体中文(澳门)',
      country: '中国澳门',
      flag: 'mo',
      prompt: 'You must write the entire article in Traditional Chinese (繁體中文). Use appropriate Chinese punctuation and expressions common in Macau.'
    },
    'zh-SG': {
      name: '简体中文(新加坡)',
      country: '新加坡',
      flag: 'sg',
      prompt: 'You must write the entire article in Simplified Chinese (简体中文). Use appropriate Chinese punctuation and expressions common in Singapore.'
    },
    'ja-JP': {
      name: '日语',
      country: '日本',
      flag: 'jp',
      prompt: 'You must write the entire article in Japanese (日本語). Use appropriate Japanese grammar, keigo, and expressions.'
    },
    'ko-KR': {
      name: '韩语',
      country: '韩国',
      flag: 'kr',
      prompt: 'You must write the entire article in Korean (한국어). Use appropriate Korean grammar, honorifics, and expressions.'
    },
    'mn-MN': {
      name: '蒙古语',
      country: '蒙古',
      flag: 'mn',
      prompt: 'You must write the entire article in Mongolian (Монгол хэл). Use appropriate Mongolian grammar, punctuation, and expressions.'
    },
    'ug-CN': {
      name: '维吾尔语',
      country: '中国',
      flag: 'cn',
      prompt: 'You must write the entire article in Uyghur (ئۇيغۇرچە). Use appropriate Uyghur grammar, punctuation, and expressions.'
    },
    'bo-CN': {
      name: '藏语',
      country: '中国',
      flag: 'cn',
      prompt: 'You must write the entire article in Tibetan (བོད་ཡིག). Use appropriate Tibetan grammar, punctuation, and expressions.'
    }
  },
  '东南亚': {
    'vi-VN': {
      name: '越南语',
      country: '越南',
      flag: 'vn',
      prompt: 'You must write the entire article in Vietnamese (Tiếng Việt). Use appropriate Vietnamese grammar, punctuation, and expressions.'
    },
    'th-TH': {
      name: '泰语',
      country: '泰国',
      flag: 'th',
      prompt: 'You must write the entire article in Thai (ภาษาไทย). Use appropriate Thai grammar, punctuation, and expressions.'
    },
    'id-ID': {
      name: '印尼语',
      country: '印度尼西亚',
      flag: 'id',
      prompt: 'You must write the entire article in Indonesian (Bahasa Indonesia). Use appropriate Indonesian grammar, punctuation, and expressions.'
    },
    'ms-MY': {
      name: '马来语',
      country: '马来西亚',
      flag: 'my',
      prompt: 'You must write the entire article in Malay (Bahasa Melayu). Use appropriate Malay grammar, punctuation, and expressions.'
    },
    'fil-PH': {
      name: '菲律宾语',
      country: '菲律宾',
      flag: 'ph',
      prompt: 'You must write the entire article in Filipino. Use appropriate Filipino grammar, punctuation, and expressions.'
    },
    'my-MM': {
      name: '缅甸语',
      country: '缅甸',
      flag: 'mm',
      prompt: 'You must write the entire article in Burmese (မြန်မာဘာသာ). Use appropriate Burmese grammar, punctuation, and expressions.'
    },
    'km-KH': {
      name: '柬埔寨语',
      country: '柬埔寨',
      flag: 'kh',
      prompt: 'You must write the entire article in Khmer (ភាសាខ្មែរ). Use appropriate Khmer grammar, punctuation, and expressions.'
    },
    'lo-LA': {
      name: '老挝语',
      country: '老挝',
      flag: 'la',
      prompt: 'You must write the entire article in Lao (ພາສາລາວ). Use appropriate Lao grammar, punctuation, and expressions.'
    },
    'ceb-PH': {
      name: '宿务语',
      country: '菲律宾',
      flag: 'ph',
      prompt: 'You must write the entire article in Cebuano (Sinugboanon). Use appropriate Cebuano grammar, punctuation, and expressions.'
    },
    'jv-ID': {
      name: '爪哇语',
      country: '印度尼西亚',
      flag: 'id',
      prompt: 'You must write the entire article in Javanese (Basa Jawa). Use appropriate Javanese grammar, punctuation, and expressions.'
    },
    'su-ID': {
      name: '巽他语',
      country: '印度尼西亚',
      flag: 'id',
      prompt: 'You must write the entire article in Sundanese (Basa Sunda). Use appropriate Sundanese grammar, punctuation, and expressions.'
    },
    'ms-BN': {
      name: '马来语(文莱)',
      country: '文莱',
      flag: 'bn',
      prompt: 'You must write the entire article in Malay (Bahasa Melayu). Use appropriate Malay grammar, punctuation, and expressions common in Brunei.'
    },
    'tet-TL': {
      name: '德顿语',
      country: '东帝汶',
      flag: 'tl',
      prompt: 'You must write the entire article in Tetum (Tetun). Use appropriate Tetum grammar, punctuation, and expressions.'
    }
  },
  '南亚': {
    'hi-IN': {
      name: '印地语',
      country: '印度',
      flag: 'in',
      prompt: 'You must write the entire article in Hindi (हिन्दी). Use appropriate Hindi grammar, punctuation, and expressions.'
    },
    'bn-IN': {
      name: '孟加拉语(印度)',
      country: '印度',
      flag: 'in',
      prompt: 'You must write the entire article in Bengali (বাংলা). Use appropriate Bengali grammar, punctuation, and expressions common in India.'
    },
    'bn-BD': {
      name: '孟加拉语',
      country: '孟加拉国',
      flag: 'bd',
      prompt: 'You must write the entire article in Bengali (বাংলা). Use appropriate Bengali grammar, punctuation, and expressions common in Bangladesh.'
    },
    'ta-IN': {
      name: '泰米尔语',
      country: '印度',
      flag: 'in',
      prompt: 'You must write the entire article in Tamil (தமிழ்). Use appropriate Tamil grammar, punctuation, and expressions.'
    },
    'te-IN': {
      name: '泰卢固语',
      country: '印度',
      flag: 'in',
      prompt: 'You must write the entire article in Telugu (తెలుగు). Use appropriate Telugu grammar, punctuation, and expressions.'
    },
    'mr-IN': {
      name: '马拉地语',
      country: '印度',
      flag: 'in',
      prompt: 'You must write the entire article in Marathi (मराठी). Use appropriate Marathi grammar, punctuation, and expressions.'
    },
    'ur-PK': {
      name: '乌尔都语',
      country: '巴基斯坦',
      flag: 'pk',
      prompt: 'You must write the entire article in Urdu (اردو). Use appropriate Urdu grammar, punctuation, and expressions.'
    },
    'si-LK': {
      name: '僧伽罗语',
      country: '斯里兰卡',
      flag: 'lk',
      prompt: 'You must write the entire article in Sinhala (සිංහල). Use appropriate Sinhala grammar, punctuation, and expressions.'
    },
    'ne-NP': {
      name: '尼泊尔语',
      country: '尼泊尔',
      flag: 'np',
      prompt: 'You must write the entire article in Nepali (नेपाली). Use appropriate Nepali grammar, punctuation, and expressions.'
    },
    'pa-IN': {
      name: '旁遮普语',
      country: '印度',
      flag: 'in',
      prompt: 'You must write the entire article in Punjabi (ਪੰਜਾਬੀ). Use appropriate Punjabi grammar, punctuation, and expressions.'
    },
    'gu-IN': {
      name: '古吉拉特语',
      country: '印度',
      flag: 'in',
      prompt: 'You must write the entire article in Gujarati (ગુજરાતી). Use appropriate Gujarati grammar, punctuation, and expressions.'
    },
    'kn-IN': {
      name: '卡纳达语',
      country: '印度',
      flag: 'in',
      prompt: 'You must write the entire article in Kannada (ಕನ್ನಡ). Use appropriate Kannada grammar, punctuation, and expressions.'
    },
    'ml-IN': {
      name: '马拉雅拉姆语',
      country: '印度',
      flag: 'in',
      prompt: 'You must write the entire article in Malayalam (മലയാളം). Use appropriate Malayalam grammar, punctuation, and expressions.'
    },
    'or-IN': {
      name: '奥里亚语',
      country: '印度',
      flag: 'in',
      prompt: 'You must write the entire article in Odia (ଓଡ଼ିଆ). Use appropriate Odia grammar, punctuation, and expressions.'
    },
    'as-IN': {
      name: '阿萨姆语',
      country: '印度',
      flag: 'in',
      prompt: 'You must write the entire article in Assamese (অসমীয়া). Use appropriate Assamese grammar, punctuation, and expressions.'
    },
    'sd-PK': {
      name: '信德语',
      country: '巴基斯坦',
      flag: 'pk',
      prompt: 'You must write the entire article in Sindhi (سنڌي). Use appropriate Sindhi grammar, punctuation, and expressions.'
    },
    'dv-MV': {
      name: '迪维希语',
      country: '马尔代夫',
      flag: 'mv',
      prompt: 'You must write the entire article in Dhivehi (ދިވެހި). Use appropriate Dhivehi grammar, punctuation, and expressions.'
    },
    'dz-BT': {
      name: '宗卡语',
      country: '不丹',
      flag: 'bt',
      prompt: 'You must write the entire article in Dzongkha (རྫོང་ཁ). Use appropriate Dzongkha grammar, punctuation, and expressions.'
    }
  },
  '中亚': {
    'kk-KZ': {
      name: '哈萨克语',
      country: '哈萨克斯坦',
      flag: 'kz',
      prompt: 'You must write the entire article in Kazakh (Қазақ тілі). Use appropriate Kazakh grammar, punctuation, and expressions.'
    },
    'ky-KG': {
      name: '吉尔吉斯语',
      country: '吉尔吉斯斯坦',
      flag: 'kg',
      prompt: 'You must write the entire article in Kyrgyz (Кыргыз тили). Use appropriate Kyrgyz grammar, punctuation, and expressions.'
    },
    'tg-TJ': {
      name: '塔吉克语',
      country: '塔吉克斯坦',
      flag: 'tj',
      prompt: 'You must write the entire article in Tajik (Тоҷикӣ). Use appropriate Tajik grammar, punctuation, and expressions.'
    },
    'tk-TM': {
      name: '土库曼语',
      country: '土库曼斯坦',
      flag: 'tm',
      prompt: 'You must write the entire article in Turkmen (Türkmen dili). Use appropriate Turkmen grammar, punctuation, and expressions.'
    },
    'uz-UZ': {
      name: '乌兹别克语',
      country: '乌兹别克斯坦',
      flag: 'uz',
      prompt: 'You must write the entire article in Uzbek (Ozbek tili). Use appropriate Uzbek grammar, punctuation, and expressions.'
    }
  },
  '欧洲': {
    'en-GB': {
      name: '英语(英国)',
      country: '英国',
      flag: 'gb',
      prompt: 'You must write the entire article in British English. Use appropriate British spelling, punctuation, and expressions.'
    },
    'fr-FR': {
      name: '法语',
      country: '法国',
      flag: 'fr',
      prompt: 'You must write the entire article in French (Francais). Use appropriate French grammar, punctuation, and expressions.'
    },
    'de-DE': {
      name: '德语',
      country: '德国',
      flag: 'de',
      prompt: 'You must write the entire article in German (Deutsch). Use appropriate German grammar, punctuation, and expressions.'
    },
    'es-ES': {
      name: '西班牙语',
      country: '西班牙',
      flag: 'es',
      prompt: 'You must write the entire article in Spanish (Espanol). Use appropriate Spanish grammar, punctuation, and expressions.'
    },
    'it-IT': {
      name: '意大利语',
      country: '意大利',
      flag: 'it',
      prompt: 'You must write the entire article in Italian (Italiano). Use appropriate Italian grammar, punctuation, and expressions.'
    },
    'ru-RU': {
      name: '俄语',
      country: '俄罗斯',
      flag: 'ru',
      prompt: 'You must write the entire article in Russian (Russkiy). Use appropriate Russian grammar, punctuation, and expressions.'
    },
    'pl-PL': {
      name: '波兰语',
      country: '波兰',
      flag: 'pl',
      prompt: 'You must write the entire article in Polish (Polski). Use appropriate Polish grammar, punctuation, and expressions.'
    },
    'uk-UA': {
      name: '乌克兰语',
      country: '乌克兰',
      flag: 'ua',
      prompt: 'You must write the entire article in Ukrainian (Українська). Use appropriate Ukrainian grammar, punctuation, and expressions.'
    },
    'nl-NL': {
      name: '荷兰语',
      country: '荷兰',
      flag: 'nl',
      prompt: 'You must write the entire article in Dutch (Nederlands). Use appropriate Dutch grammar, punctuation, and expressions.'
    },
    'pt-PT': {
      name: '葡萄牙语',
      country: '葡萄牙',
      flag: 'pt',
      prompt: 'You must write the entire article in Portuguese (Português). Use appropriate Portuguese grammar, punctuation, and expressions.'
    },
    'sv-SE': {
      name: '瑞典语',
      country: '瑞典',
      flag: 'se',
      prompt: 'You must write the entire article in Swedish (Svenska). Use appropriate Swedish grammar, punctuation, and expressions.'
    },
    'da-DK': {
      name: '丹麦语',
      country: '丹麦',
      flag: 'dk',
      prompt: 'You must write the entire article in Danish (Dansk). Use appropriate Danish grammar, punctuation, and expressions.'
    },
    'fi-FI': {
      name: '芬兰语',
      country: '芬兰',
      flag: 'fi',
      prompt: 'You must write the entire article in Finnish (Suomi). Use appropriate Finnish grammar, punctuation, and expressions.'
    },
    'el-GR': {
      name: '希腊语',
      country: '希腊',
      flag: 'gr',
      prompt: 'You must write the entire article in Greek (Ελληνικά). Use appropriate Greek grammar, punctuation, and expressions.'
    },
    'hu-HU': {
      name: '匈牙利语',
      country: '匈牙利',
      flag: 'hu',
      prompt: 'You must write the entire article in Hungarian (Magyar). Use appropriate Hungarian grammar, punctuation, and expressions.'
    },
    'cs-CZ': {
      name: '捷克语',
      country: '捷克',
      flag: 'cz',
      prompt: 'You must write the entire article in Czech (Čeština). Use appropriate Czech grammar, punctuation, and expressions.'
    },
    'ro-RO': {
      name: '罗马尼亚语',
      country: '罗马尼亚',
      flag: 'ro',
      prompt: 'You must write the entire article in Romanian (Română). Use appropriate Romanian grammar, punctuation, and expressions.'
    },
    'bg-BG': {
      name: '保加利亚语',
      country: '保加利亚',
      flag: 'bg',
      prompt: 'You must write the entire article in Bulgarian (Български). Use appropriate Bulgarian grammar, punctuation, and expressions.'
    },
    'sk-SK': {
      name: '斯洛伐克语',
      country: '斯洛伐克',
      flag: 'sk',
      prompt: 'You must write the entire article in Slovak (Slovenčina). Use appropriate Slovak grammar, punctuation, and expressions.'
    },
    'hr-HR': {
      name: '克罗地亚语',
      country: '克罗地亚',
      flag: 'hr',
      prompt: 'You must write the entire article in Croatian (Hrvatski). Use appropriate Croatian grammar, punctuation, and expressions.'
    },
    'sr-RS': {
      name: '塞尔维亚语',
      country: '塞尔维亚',
      flag: 'rs',
      prompt: 'You must write the entire article in Serbian (Српски). Use appropriate Serbian grammar, punctuation, and expressions.'
    },
    'lt-LT': {
      name: '立陶宛语',
      country: '立陶宛',
      flag: 'lt',
      prompt: 'You must write the entire article in Lithuanian (Lietuvių). Use appropriate Lithuanian grammar, punctuation, and expressions.'
    },
    'lv-LV': {
      name: '拉脱维亚语',
      country: '拉脱维亚',
      flag: 'lv',
      prompt: 'You must write the entire article in Latvian (Latviešu). Use appropriate Latvian grammar, punctuation, and expressions.'
    },
    'et-EE': {
      name: '爱沙尼亚语',
      country: '爱沙尼亚',
      flag: 'ee',
      prompt: 'You must write the entire article in Estonian (Eesti). Use appropriate Estonian grammar, punctuation, and expressions.'
    },
    'no-NO': {
      name: '挪威语',
      country: '挪威',
      flag: 'no',
      prompt: 'You must write the entire article in Norwegian (Norsk). Use appropriate Norwegian grammar, punctuation, and expressions.'
    },
    'is-IS': {
      name: '冰岛语',
      country: '冰岛',
      flag: 'is',
      prompt: 'You must write the entire article in Icelandic (Íslenska). Use appropriate Icelandic grammar, punctuation, and expressions.'
    },
    'ga-IE': {
      name: '爱尔兰语',
      country: '爱尔兰',
      flag: 'ie',
      prompt: 'You must write the entire article in Irish (Gaeilge). Use appropriate Irish grammar, punctuation, and expressions.'
    },
    'mt-MT': {
      name: '马耳他语',
      country: '马耳他',
      flag: 'mt',
      prompt: 'You must write the entire article in Maltese (Malti). Use appropriate Maltese grammar, punctuation, and expressions.'
    },
    'eu-ES': {
      name: '巴斯克语',
      country: '西班牙',
      flag: 'es',
      prompt: 'You must write the entire article in Basque (Euskara). Use appropriate Basque grammar, punctuation, and expressions.'
    },
    'ca-ES': {
      name: '加泰罗尼亚语',
      country: '西班牙',
      flag: 'es',
      prompt: 'You must write the entire article in Catalan (Català). Use appropriate Catalan grammar, punctuation, and expressions.'
    },
    'cy-GB': {
      name: '威尔士语',
      country: '英国',
      flag: 'gb',
      prompt: 'You must write the entire article in Welsh (Cymraeg). Use appropriate Welsh grammar, punctuation, and expressions.'
    },
    'fo-FO': {
      name: '法罗语',
      country: '法罗群岛',
      flag: 'fo',
      prompt: 'You must write the entire article in Faroese (Føroyskt). Use appropriate Faroese grammar, punctuation, and expressions.'
    },
    'kl-GL': {
      name: '格陵兰语',
      country: '格陵兰',
      flag: 'gl',
      prompt: 'You must write the entire article in Greenlandic (Kalaallisut). Use appropriate Greenlandic grammar, punctuation, and expressions.'
    },
    'gd-GB': {
      name: '苏格兰盖尔语',
      country: '英国',
      flag: 'gb',
      prompt: 'You must write the entire article in Scottish Gaelic (Gàidhlig). Use appropriate Gaelic grammar, punctuation, and expressions.'
    },
    'be-BY': {
      name: '白俄罗斯语',
      country: '白俄罗斯',
      flag: 'by',
      prompt: 'You must write the entire article in Belarusian (Беларуская). Use appropriate Belarusian grammar, punctuation, and expressions.'
    },
    'mk-MK': {
      name: '马其顿语',
      country: '北马其顿',
      flag: 'mk',
      prompt: 'You must write the entire article in Macedonian (Македонски). Use appropriate Macedonian grammar, punctuation, and expressions.'
    },
    'sq-AL': {
      name: '阿尔巴尼亚语',
      country: '阿尔巴尼亚',
      flag: 'al',
      prompt: 'You must write the entire article in Albanian (Shqip). Use appropriate Albanian grammar, punctuation, and expressions.'
    },
    'fr-MC': {
      name: '法语(摩纳哥)',
      country: '摩纳哥',
      flag: 'mc',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Monaco.'
    },
    'de-LI': {
      name: '德语(列支敦士登)',
      country: '列支敦士登',
      flag: 'li',
      prompt: 'You must write the entire article in German (Deutsch). Use appropriate German grammar, punctuation, and expressions common in Liechtenstein.'
    },
    'ca-AD': {
      name: '加泰罗尼亚语(安道尔)',
      country: '安道尔',
      flag: 'ad',
      prompt: 'You must write the entire article in Catalan (Català). Use appropriate Catalan grammar, punctuation, and expressions common in Andorra.'
    },
    'it-SM': {
      name: '意大利语(圣马力诺)',
      country: '圣马力诺',
      flag: 'sm',
      prompt: 'You must write the entire article in Italian (Italiano). Use appropriate Italian grammar, punctuation, and expressions common in San Marino.'
    },
    'it-VA': {
      name: '意大利语(梵蒂冈)',
      country: '梵蒂冈',
      flag: 'va',
      prompt: 'You must write the entire article in Italian (Italiano). Use appropriate Italian grammar, punctuation, and expressions common in Vatican City.'
    },
    'lb-LU': {
      name: '卢森堡语',
      country: '卢森堡',
      flag: 'lu',
      prompt: 'You must write the entire article in Luxembourgish (Lëtzebuergesch). Use appropriate Luxembourgish grammar, punctuation, and expressions.'
    },
    'ro-MD': {
      name: '罗马尼亚语(摩尔多瓦)',
      country: '摩尔多瓦',
      flag: 'md',
      prompt: 'You must write the entire article in Romanian (Română). Use appropriate Romanian grammar, punctuation, and expressions common in Moldova.'
    }
  },
  '美洲': {
    'en-US': {
      name: '英语(美国)',
      country: '美国',
      flag: 'us',
      prompt: 'You must write the entire article in American English. Use appropriate American spelling, punctuation, and expressions.'
    },
    'es-MX': {
      name: '西班牙语(墨西哥)',
      country: '墨西哥',
      flag: 'mx',
      prompt: 'You must write the entire article in Mexican Spanish (Espanol mexicano). Use appropriate Mexican Spanish grammar, punctuation, and expressions.'
    },
    'pt-BR': {
      name: '葡萄牙语(巴西)',
      country: '巴西',
      flag: 'br',
      prompt: 'You must write the entire article in Brazilian Portuguese (Portugues do Brasil). Use appropriate Brazilian Portuguese grammar, punctuation, and expressions.'
    },
    'fr-CA': {
      name: '法语(加拿大)',
      country: '加拿大',
      flag: 'ca',
      prompt: 'You must write the entire article in Canadian French (Français du Canada). Use appropriate Canadian French grammar, punctuation, and expressions.'
    },
    'es-AR': {
      name: '西班牙语(阿根廷)',
      country: '阿根廷',
      flag: 'ar',
      prompt: 'You must write the entire article in Argentine Spanish (Español de Argentina). Use appropriate Argentine Spanish grammar, punctuation, and expressions.'
    },
    'es-CL': {
      name: '西班牙语(智利)',
      country: '智利',
      flag: 'cl',
      prompt: 'You must write the entire article in Chilean Spanish (Español de Chile). Use appropriate Chilean Spanish grammar, punctuation, and expressions.'
    },
    'es-CO': {
      name: '西班牙语(哥伦比亚)',
      country: '哥伦比亚',
      flag: 'co',
      prompt: 'You must write the entire article in Colombian Spanish (Español de Colombia). Use appropriate Colombian Spanish grammar, punctuation, and expressions.'
    },
    'es-PE': {
      name: '西班牙语(秘鲁)',
      country: '秘鲁',
      flag: 'pe',
      prompt: 'You must write the entire article in Peruvian Spanish (Español peruano). Use appropriate Peruvian Spanish grammar, punctuation, and expressions.'
    },
    'es-VE': {
      name: '西班牙语(委内瑞拉)',
      country: '委内瑞拉',
      flag: 've',
      prompt: 'You must write the entire article in Venezuelan Spanish (Español venezolano). Use appropriate Venezuelan Spanish grammar, punctuation, and expressions.'
    },
    'qu-PE': {
      name: '克丘亚语',
      country: '秘鲁',
      flag: 'pe',
      prompt: 'You must write the entire article in Quechua (Runa Simi). Use appropriate Quechua grammar, punctuation, and expressions.'
    },
    'ay-BO': {
      name: '艾马拉语',
      country: '玻利维亚',
      flag: 'bo',
      prompt: 'You must write the entire article in Aymara (Aymar aru). Use appropriate Aymara grammar, punctuation, and expressions.'
    },
    'gn-PY': {
      name: '瓜拉尼语',
      country: '巴拉圭',
      flag: 'py',
      prompt: 'You must write the entire article in Guarani (Guarani). Use appropriate Guarani grammar, punctuation, and expressions.'
    },
    'en-AG': {
      name: '英语(安提瓜和巴布达)',
      country: '安提瓜和巴布达',
      flag: 'ag',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Antigua and Barbuda.'
    },
    'en-BS': {
      name: '英语(巴哈马)',
      country: '巴哈马',
      flag: 'bs',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in The Bahamas.'
    },
    'en-BB': {
      name: '英语(巴巴多斯)',
      country: '巴巴多斯',
      flag: 'bb',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Barbados.'
    },
    'en-BZ': {
      name: '英语(伯利兹)',
      country: '伯利兹',
      flag: 'bz',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Belize.'
    },
    'en-DM': {
      name: '英语(多米尼克)',
      country: '多米尼克',
      flag: 'dm',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Dominica.'
    },
    'es-DO': {
      name: '西班牙语(多米尼加)',
      country: '多米尼加共和国',
      flag: 'do',
      prompt: 'You must write the entire article in Spanish (Español). Use appropriate Spanish grammar, punctuation, and expressions common in Dominican Republic.'
    },
    'en-GD': {
      name: '英语(格林纳达)',
      country: '格林纳达',
      flag: 'gd',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Grenada.'
    },
    'en-GY': {
      name: '英语(圭亚那)',
      country: '圭亚那',
      flag: 'gy',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Guyana.'
    },
    'ht-HT': {
      name: '海地克里奥尔语',
      country: '海地',
      flag: 'ht',
    }
  },
  '大洋洲': {
    'en-AU': {
      name: '英语(澳大利亚)',
      country: '澳大利亚',
      flag: 'au',
      prompt: 'You must write the entire article in Australian English. Use appropriate Australian spelling, punctuation, and expressions.'
    },
    'en-NZ': {
      name: '英语(新西兰)',
      country: '新西兰',
      flag: 'nz',
      prompt: 'You must write the entire article in New Zealand English. Use appropriate New Zealand spelling, punctuation, and expressions.'
    },
    'en-KI': {
      name: '英语(基里巴斯)',
      country: '基里巴斯',
      flag: 'ki',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Kiribati.'
    },
    'mh-MH': {
      name: '马绍尔语',
      country: '马绍尔群岛',
      flag: 'mh',
      prompt: 'You must write the entire article in Marshallese (Kajin Majel). Use appropriate Marshallese grammar, punctuation, and expressions.'
    },
    'en-FM': {
      name: '英语(密克罗尼西亚)',
      country: '密克罗尼西亚',
      flag: 'fm',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Micronesia.'
    },
    'na-NR': {
      name: '瑙鲁语',
      country: '瑙鲁',
      flag: 'nr',
      prompt: 'You must write the entire article in Nauruan (Dorerin Naoero). Use appropriate Nauruan grammar, punctuation, and expressions.'
    },
    'pau-PW': {
      name: '帕劳语',
      country: '帕劳',
      flag: 'pw',
      prompt: 'You must write the entire article in Palauan (a tekoi er a Belau). Use appropriate Palauan grammar, punctuation, and expressions.'
    },
    'en-PG': {
      name: '英语(巴布亚新几内亚)',
      country: '巴布亚新几内亚',
      flag: 'pg',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Papua New Guinea.'
    },
    'en-SB': {
      name: '英语(所罗门群岛)',
      country: '所罗门群岛',
      flag: 'sb',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Solomon Islands.'
    },
    'tvl-TV': {
      name: '图瓦卢语',
      country: '图瓦卢',
      flag: 'tv',
      prompt: 'You must write the entire article in Tuvaluan (Te Ggana Tuuvalu). Use appropriate Tuvaluan grammar, punctuation, and expressions.'
    },
    'bi-VU': {
      name: '比斯拉马语',
      country: '瓦努阿图',
      flag: 'vu',
      prompt: 'You must write the entire article in Bislama. Use appropriate Bislama grammar, punctuation, and expressions.'
    },
    'sm-WS': {
      name: '萨摩亚语',
      country: '萨摩亚',
      flag: 'ws',
      prompt: 'You must write the entire article in Samoan (Gagana Samoa). Use appropriate Samoan grammar, punctuation, and expressions.'
    },
    'to-TO': {
      name: '汤加语',
      country: '汤加',
      flag: 'to',
      prompt: 'You must write the entire article in Tongan (Faka Tonga). Use appropriate Tongan grammar, punctuation, and expressions.'
    },
    'fj-FJ': {
      name: '斐济语',
      country: '斐济',
      flag: 'fj',
      prompt: 'You must write the entire article in Fijian (Na Vosa Vakaviti). Use appropriate Fijian grammar, punctuation, and expressions.'
    },
    'ty-PF': {
      name: '塔希提语',
      country: '法属波利尼西亚',
      flag: 'pf',
      prompt: 'You must write the entire article in Tahitian (Reo Mātauranga). Use appropriate Tahitian grammar, punctuation, and expressions.'
    }
  },
  '西亚': {
    'ar-SA': {
      name: '阿拉伯语(沙特)',
      country: '沙特阿拉伯',
      flag: 'sa',
      prompt: 'You must write the entire article in Arabic. Use appropriate Arabic grammar, punctuation, and expressions common in Saudi Arabia.'
    },
    'ar-AE': {
      name: '阿拉伯语(阿联酋)',
      country: '阿联酋',
      flag: 'ae',
      prompt: 'You must write the entire article in Arabic (العربية). Use appropriate Arabic grammar, punctuation, and expressions common in UAE.'
    },
    'fa-IR': {
      name: '波斯语',
      country: '伊朗',
      flag: 'ir',
      prompt: 'You must write the entire article in Persian (Farsi). Use appropriate Persian grammar, punctuation, and expressions.'
    },
    'tr-TR': {
      name: '土耳其语',
      country: '土耳其',
      flag: 'tr',
      prompt: 'You must write the entire article in Turkish (Turkce). Use appropriate Turkish grammar, punctuation, and expressions.'
    },
    'he-IL': {
      name: '希伯来语',
      country: '以色列',
      flag: 'il',
      prompt: 'You must write the entire article in Hebrew (עברית). Use appropriate Hebrew grammar, punctuation, and expressions.'
    },
    'ar-IQ': {
      name: '阿拉伯语(伊拉克)',
      country: '伊拉克',
      flag: 'iq',
      prompt: 'You must write the entire article in Iraqi Arabic (اللهجة العراقية). Use appropriate Iraqi Arabic grammar, punctuation, and expressions.'
    },
    'ar-JO': {
      name: '阿拉伯语(约旦)',
      country: '约旦',
      flag: 'jo',
      prompt: 'You must write the entire article in Jordanian Arabic (اللهجة الأردنية). Use appropriate Jordanian Arabic grammar, punctuation, and expressions.'
    },
    'ar-SY': {
      name: '阿拉伯语(叙利亚)',
      country: '叙利亚',
      flag: 'sy',
      prompt: 'You must write the entire article in Arabic (العربية). Use appropriate Arabic grammar, punctuation, and expressions common in Syria.'
    },
    'ar-LB': {
      name: '阿拉伯语(黎巴嫩)',
      country: '黎巴嫩',
      flag: 'lb',
      prompt: 'You must write the entire article in Arabic (العربية). Use appropriate Arabic grammar, punctuation, and expressions common in Lebanon.'
    },
    'ar-KW': {
      name: '阿拉伯语(科威特)',
      country: '科威特',
      flag: 'kw',
      prompt: 'You must write the entire article in Arabic (العربية). Use appropriate Arabic grammar, punctuation, and expressions common in Kuwait.'
    },
    'ar-OM': {
      name: '阿拉伯语(阿曼)',
      country: '阿曼',
      flag: 'om',
      prompt: 'You must write the entire article in Arabic (العربية). Use appropriate Arabic grammar, punctuation, and expressions common in Oman.'
    },
    'ar-BH': {
      name: '阿拉伯语(巴林)',
      country: '巴林',
      flag: 'bh',
      prompt: 'You must write the entire article in Arabic (العربية). Use appropriate Arabic grammar, punctuation, and expressions common in Bahrain.'
    },
    'ar-QA': {
      name: '阿拉伯语(卡塔尔)',
      country: '卡塔尔',
      flag: 'qa',
      prompt: 'You must write the entire article in Arabic (العربية). Use appropriate Arabic grammar, punctuation, and expressions common in Qatar.'
    },
    'hy-AM': {
      name: '亚美尼亚语',
      country: '亚美尼亚',
      flag: 'am',
      prompt: 'You must write the entire article in Armenian (Հայերեն). Use appropriate Armenian grammar, punctuation, and expressions.'
    },
    'ka-GE': {
      name: '格鲁吉亚语',
      country: '格鲁吉亚',
      flag: 'ge',
      prompt: 'You must write the entire article in Georgian (ქართული). Use appropriate Georgian grammar, punctuation, and expressions.'
    },
    'az-AZ': {
      name: '阿塞拜疆语',
      country: '阿塞拜疆',
      flag: 'az',
      prompt: 'You must write the entire article in Azerbaijani (Azərbaycan dili). Use appropriate Azerbaijani grammar, punctuation, and expressions.'
    }
  },
  '非洲': {
    'ar-EG': {
      name: '阿拉伯语(埃及)',
      country: '埃及',
      flag: 'eg',
      prompt: 'You must write the entire article in Arabic. Use appropriate Arabic grammar, punctuation, and expressions common in Egypt.'
    },
    'am-ET': {
      name: '阿姆哈拉语',
      country: '埃塞俄比亚',
      flag: 'et',
      prompt: 'You must write the entire article in Amharic. Use appropriate Amharic grammar, punctuation, and expressions.'
    },
    'sw-KE': {
      name: '斯瓦希里语',
      country: '肯尼亚',
      flag: 'ke',
      prompt: 'You must write the entire article in Swahili (Kiswahili). Use appropriate Swahili grammar, punctuation, and expressions.'
    },
    'zu-ZA': {
      name: '祖鲁语',
      country: '南非',
      flag: 'za',
      prompt: 'You must write the entire article in Zulu (isiZulu). Use appropriate Zulu grammar, punctuation, and expressions.'
    },
    'ar-MA': {
      name: '阿拉伯语(摩洛哥)',
      country: '摩洛哥',
      flag: 'ma',
      prompt: 'You must write the entire article in Moroccan Arabic (الدارجة المغربية). Use appropriate Moroccan Arabic grammar, punctuation, and expressions.'
    },
    'ar-TN': {
      name: '阿拉伯语(突尼斯)',
      country: '突尼斯',
      flag: 'tn',
      prompt: 'You must write the entire article in Tunisian Arabic (الدارجة التونسية). Use appropriate Tunisian Arabic grammar, punctuation, and expressions.'
    },
    'ha-NG': {
      name: '豪萨语',
      country: '尼日利亚',
      flag: 'ng',
      prompt: 'You must write the entire article in Hausa (Harshen Hausa). Use appropriate Hausa grammar, punctuation, and expressions.'
    },
    'yo-NG': {
      name: '约鲁巴语',
      country: '尼日利亚',
      flag: 'ng',
      prompt: 'You must write the entire article in Yoruba (Èdè Yorùbá). Use appropriate Yoruba grammar, punctuation, and expressions.'
    },
    'ar-DZ': {
      name: '阿拉伯语(阿尔及利亚)',
      country: '阿尔及利亚',
      flag: 'dz',
      prompt: 'You must write the entire article in Algerian Arabic (العربية الجزائرية). Use appropriate Algerian Arabic grammar, punctuation, and expressions.'
    },
    'ar-LY': {
      name: '阿拉伯语(利比亚)',
      country: '利比亚',
      flag: 'ly',
      prompt: 'You must write the entire article in Libyan Arabic (العربية الليبية). Use appropriate Libyan Arabic grammar, punctuation, and expressions.'
    },
    'ar-SD': {
      name: '阿拉伯语(苏丹)',
      country: '苏丹',
      flag: 'sd',
      prompt: 'You must write the entire article in Sudanese Arabic (العربية السودانية). Use appropriate Sudanese Arabic grammar, punctuation, and expressions.'
    },
    'ig-NG': {
      name: '伊博语',
      country: '尼日利亚',
      flag: 'ng',
      prompt: 'You must write the entire article in Igbo (Igbo). Use appropriate Igbo grammar, punctuation, and expressions.'
    },
    'wo-SN': {
      name: '沃洛夫语',
      country: '塞内加尔',
      flag: 'sn',
      prompt: 'You must write the entire article in Wolof (Wolof). Use appropriate Wolof grammar, punctuation, and expressions.'
    },
    'xh-ZA': {
      name: '科萨语',
      country: '南非',
      flag: 'za',
      prompt: 'You must write the entire article in Xhosa (isiXhosa). Use appropriate Xhosa grammar, punctuation, and expressions.'
    },
    'af-ZA': {
      name: '南非荷兰语',
      country: '南非',
      flag: 'za',
      prompt: 'You must write the entire article in Afrikaans (Afrikaans). Use appropriate Afrikaans grammar, punctuation, and expressions.'
    },
    'sn-ZW': {
      name: '绍纳语',
      country: '津巴布韦',
      flag: 'zw',
      prompt: 'You must write the entire article in Shona (ChiShona). Use appropriate Shona grammar, punctuation, and expressions.'
    },
    'rw-RW': {
      name: '卢旺达语',
      country: '卢旺达',
      flag: 'rw',
      prompt: 'You must write the entire article in Kinyarwanda (Ikinyarwanda). Use appropriate Kinyarwanda grammar, punctuation, and expressions.'
    },
    'ny-MW': {
      name: '齐切瓦语',
      country: '马拉维',
      flag: 'mw',
      prompt: 'You must write the entire article in Chichewa (ChiCheŵa). Use appropriate Chichewa grammar, punctuation, and expressions.'
    },
    'mg-MG': {
      name: '马达加斯加语',
      country: '马达加斯加',
      flag: 'mg',
      prompt: 'You must write the entire article in Malagasy (Malagasy). Use appropriate Malagasy grammar, punctuation, and expressions.'
    },
    'so-SO': {
      name: '索马里语(索马里)',
      country: '索马里',
      flag: 'so',
      prompt: 'You must write the entire article in Somali (Soomaali). Use appropriate Somali grammar, punctuation, and expressions common in Somalia.'
    },
    'so-DJ': {
      name: '索马里语(吉布提)',
      country: '吉布提',
      flag: 'dj',
      prompt: 'You must write the entire article in Somali (Soomaali). Use appropriate Somali grammar, punctuation, and expressions common in Djibouti.'
    },
    'so-ET': {
      name: '索马里语(埃塞俄比亚)',
      country: '埃塞俄比亚',
      flag: 'et',
      prompt: 'You must write the entire article in Somali (Soomaali). Use appropriate Somali grammar, punctuation, and expressions common in Ethiopia.'
    },
    'so-KE': {
      name: '索马里语(肯尼亚)',
      country: '肯尼亚',
      flag: 'ke',
      prompt: 'You must write the entire article in Somali (Soomaali). Use appropriate Somali grammar, punctuation, and expressions common in Kenya.'
    },
    'fr-BJ': {
      name: '法语(贝宁)',
      country: '贝宁',
      flag: 'bj',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Benin.'
    },
    'fr-BF': {
      name: '法语(布基纳法索)',
      country: '布基纳法索',
      flag: 'bf',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Burkina Faso.'
    },
    'rn-BI': {
      name: '科尔尼语',
      country: '布隆迪',
      flag: 'bi',
      prompt: 'You must write the entire article in Kirundi (Ikirundi). Use appropriate Kirundi grammar, punctuation, and expressions.'
    },
    'fr-CM': {
      name: '法语(喀麦隆)',
      country: '喀麦隆',
      flag: 'cm',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Cameroon.'
    },
    'fr-CF': {
      name: '法语(中非)',
      country: '中非共和国',
      flag: 'cf',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Central African Republic.'
    },
    'fr-TD': {
      name: '法语(乍得)',
      country: '乍得',
      flag: 'td',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Chad.'
    },
    'fr-KM': {
      name: '法语(科摩罗)',
      country: '科摩罗',
      flag: 'km',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Comoros.'
    },
    'fr-CG': {
      name: '法语(刚果共和国)',
      country: '刚果共和国',
      flag: 'cg',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Republic of the Congo.'
    },
    'fr-CD': {
      name: '法语(刚果民主共和国)',
      country: '刚果民主共和国',
      flag: 'cd',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Democratic Republic of the Congo.'
    },
    'fr-CI': {
      name: '法语(科特迪瓦)',
      country: '科特迪瓦',
      flag: 'ci',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Ivory Coast.'
    },
    'es-GQ': {
      name: '西班牙语(赤道几内亚)',
      country: '赤道几内亚',
      flag: 'gq',
      prompt: 'You must write the entire article in Spanish (Español). Use appropriate Spanish grammar, punctuation, and expressions common in Equatorial Guinea.'
    },
    'ti-ER': {
      name: '提格雷尼亚语',
      country: '厄立特里亚',
      flag: 'er',
      prompt: 'You must write the entire article in Tigrinya (ትግርኛ). Use appropriate Tigrinya grammar, punctuation, and expressions.'
    },
    'fr-GA': {
      name: '法语(加蓬)',
      country: '加蓬',
      flag: 'ga',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Gabon.'
    },
    'en-GM': {
      name: '英语(冈比亚)',
      country: '冈比亚',
      flag: 'gm',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in The Gambia.'
    },
    'fr-GN': {
      name: '法语(几内亚)',
      country: '几内亚',
      flag: 'gn',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Guinea.'
    },
    'pt-GW': {
      name: '葡萄牙语(几内亚比绍)',
      country: '几内亚比绍',
      flag: 'gw',
      prompt: 'You must write the entire article in Portuguese (Português). Use appropriate Portuguese grammar, punctuation, and expressions common in Guinea-Bissau.'
    },
    'st-LS': {
      name: '索托语',
      country: '莱索托',
      flag: 'ls',
      prompt: 'You must write the entire article in Sesotho (Sesotho). Use appropriate Sesotho grammar, punctuation, and expressions.'
    },
    'en-LR': {
      name: '英语(利比里亚)',
      country: '利比里亚',
      flag: 'lr',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Liberia.'
    },
    'fr-ML': {
      name: '法语(马里)',
      country: '马里',
      flag: 'ml',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Mali.'
    },
    'ar-MR': {
      name: '阿拉伯语(毛里塔尼亚)',
      country: '毛里塔尼亚',
      flag: 'mr',
      prompt: 'You must write the entire article in Arabic (العربية). Use appropriate Arabic grammar, punctuation, and expressions common in Mauritania.'
    },
    'en-MU': {
      name: '英语(毛里求斯)',
      country: '毛里求斯',
      flag: 'mu',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Mauritius.'
    },
    'pt-MZ': {
      name: '葡萄牙语(莫桑比克)',
      country: '莫桑比克',
      flag: 'mz',
      prompt: 'You must write the entire article in Portuguese (Português). Use appropriate Portuguese grammar, punctuation, and expressions common in Mozambique.'
    },
    'en-NA': {
      name: '英语(纳米比亚)',
      country: '纳米比亚',
      flag: 'na',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Namibia.'
    },
    'fr-NE': {
      name: '法语(尼日尔)',
      country: '尼日尔',
      flag: 'ne',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Niger.'
    },
    'pt-ST': {
      name: '葡萄牙语(圣多美和普林西比)',
      country: '圣多美和普林西比',
      flag: 'st',
      prompt: 'You must write the entire article in Portuguese (Português). Use appropriate Portuguese grammar, punctuation, and expressions common in São Tomé and Príncipe.'
    },
    'fr-SC': {
      name: '法语(塞舌尔)',
      country: '塞舌尔',
      flag: 'sc',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Seychelles.'
    },
    'en-SL': {
      name: '英语(塞拉利昂)',
      country: '塞拉利昂',
      flag: 'sl',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Sierra Leone.'
    },
    'fr-TG': {
      name: '法语(多哥)',
      country: '多哥',
      flag: 'tg',
      prompt: 'You must write the entire article in French (Français). Use appropriate French grammar, punctuation, and expressions common in Togo.'
    },
    'en-ZM': {
      name: '英语(赞比亚)',
      country: '赞比亚',
      flag: 'zm',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Zambia.'
    },
    'en-BW': {
      name: '英语(博茨瓦纳)',
      country: '博茨瓦纳',
      flag: 'bw',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Botswana.'
    },
    'ar-KM': {
      name: '阿拉伯语(科摩罗)',
      country: '科摩罗',
      flag: 'km',
      prompt: 'You must write the entire article in Arabic (العربية). Use appropriate Arabic grammar, punctuation, and expressions common in Comoros.'
    },
    'pt-CV': {
      name: '葡萄牙语(佛得角)',
      country: '佛得角',
      flag: 'cv',
      prompt: 'You must write the entire article in Portuguese (Português). Use appropriate Portuguese grammar, punctuation, and expressions common in Cape Verde.'
    },
    'en-UG': {
      name: '英语(乌干达)',
      country: '乌干达',
      flag: 'ug',
      prompt: 'You must write the entire article in English. Use appropriate English grammar, punctuation, and expressions common in Uganda.'
    }
  }
}

// 获取所有区域
export const getRegions = () => Object.keys(languageGroups)

// 获取指定区域的所有语言
export const getLanguagesByRegion = (region: string) => languageGroups[region] || {}

// 根据语言代码获取语言信息
export const getLanguageByCode = (code: string): Language | null => {
  for (const region in languageGroups) {
    if (languageGroups[region][code]) {
      return languageGroups[region][code]
    }
  }
  return null
}

// 搜索语言
export const searchLanguages = (query: string): LanguageGroups => {
  if (!query) return languageGroups

  const result: LanguageGroups = {}
  const lowerQuery = query.toLowerCase()

  for (const region in languageGroups) {
    const matchedLanguages: LanguageGroup = {}
    
    for (const code in languageGroups[region]) {
      const language = languageGroups[region][code]
      if (
        language.name.toLowerCase().includes(lowerQuery) ||
        language.country.toLowerCase().includes(lowerQuery) ||
        code.toLowerCase().includes(lowerQuery)
      ) {
        matchedLanguages[code] = language
      }
    }

    if (Object.keys(matchedLanguages).length > 0) {
      result[region] = matchedLanguages
    }
  }

  return result
} 