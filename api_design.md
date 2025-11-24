# Go Web API Design

## 功能概述

将原有Python程序重构为Go Web应用，提供RESTful API供前端调用，去除剪贴板等与操作系统相关的功能。

## API 接口设计

### 1. 获取角色列表
```
GET /api/characters
Response:
[
  {
    "id": "ema",
    "name": "樱羽艾玛",
    "emotion_count": 8,
    "font": "font3.ttf"
  },
  ...
]
```

### 2. 获取默认角色
```
GET /api/characters/current
Response:
{
  "id": "sherri",
  "name": "橘雪莉",
  "emotion_count": 7,
  "font": "font3.ttf"
}
```

### 3. 生成图片
```
POST /api/generate
Request:
{
  "type": "text",                 // 类型：目前只支持text
  "content": "示例文本内容",       // 文本内容
  "textInput": "输入的文本内容",    // 用户输入的文本（冗余字段，与content相同）
  "characterId": "sherri",        // 角色ID（可选，默认为配置文件中的默认角色，可设置为"random"表示随机）
  "emotionIndex": 1,              // 表情索引（可选，默认随机）
  "backgroundIndex": 1            // 背景索引（可选，默认随机）
}

Response:
{
  "success": true,
  "imageUrl": "/images/result_xxxxxx.png"
}
```

### 4. 获取角色表情列表
```
GET /api/characters/{characterId}/emotions
Response:
[
  {
    "id": 1,
    "name": "表情1"
  },
  ...
]
```

### 5. 获取背景列表
```
GET /api/backgrounds
Response:
[
  {
    "id": 1,
    "name": "背景1"
  },
  ...
]
```

## 无状态设计说明

后端API采用无状态设计，不保存用户选择的状态信息。所有需要的参数都通过API请求传递：

1. 默认角色从配置文件中读取
2. 如果需要随机角色，将characterId设置为"random"
3. 所有选择都在generate接口中通过参数传递

## 配置文件说明

项目使用JSON格式的配置文件来管理各种设置：

1. `config/app.json` - 应用基本配置，包括文本框坐标和默认角色
2. `config/characters.json` - 角色列表配置
3. `config/backgrounds.json` - 背景列表配置

这种设计使项目更加灵活，便于维护和扩展。