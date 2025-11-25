# 魔法少女裁判文本框生成器 API 文档

## 功能概述

本项目是一个基于Go语言实现的Web应用，提供RESTful API供前端调用，用于生成魔法少女裁判风格的文本框图片。

## API 接口设计

### 1. 获取角色列表
```
GET /api/characters

响应示例:
[
  {
    "id": "char0",
    "name": "樱羽艾玛"
  },
  {
    "id": "char1",
    "name": "二阶堂希罗"
  }
]
```

### 2. 获取默认角色
```
GET /api/characters/current

响应示例:
{
  "id": "char2",
  "name": "橘雪莉"
}
```

### 3. 生成图片
```
POST /api/generate

请求体示例:
{
  "type": "text",
  "content": "示例文本内容",
  "textInput": "输入的文本内容",
  "characterId": "char2",        // 角色ID（可选，默认为配置文件中的默认角色，可设置为"random"表示随机）
  "emotionIndex": 1,              // 表情索引（可选，默认随机）
  "backgroundIndex": 1            // 背景索引（可选，默认随机）
}

响应示例:
{
  "success": true,
  "imageData": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
  "character": "char2"
}
```

### 4. 获取角色表情列表
```
GET /api/characters/{characterId}/emotions

响应示例:
[
  {
    "id": 1,
    "name": "表情1"
  },
  {
    "id": 2,
    "name": "表情2"
  }
]
```

### 5. 获取背景列表
```
GET /api/backgrounds

响应示例:
[
  {
    "name": "背景1",
    "filename": "backgrounds/bg1.png"
  },
  {
    "name": "背景2",
    "filename": "backgrounds/bg2.png"
  }
]
```

## 无状态设计说明

后端API采用无状态设计，不保存用户选择的状态信息。所有需要的参数都通过API请求传递：

1. 默认角色从配置文件中读取
2. 如果需要随机角色，将characterId设置为"random"
3. 所有选择都在generate接口中通过参数传递

## 配置文件说明

项目使用JSON格式的配置文件来管理各种设置：

1. `config/app.json` - 应用基本配置，包括文本框坐标、默认角色和端口号
2. `config/characters.json` - 角色列表配置
3. `config/backgrounds.json` - 背景列表配置

这种设计使项目更加灵活，便于维护和扩展。