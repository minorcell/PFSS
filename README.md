# File-Upload-Demo

一个基于Golang和Gin框架实现的文件上传和管理演示项目。

## 功能特点

- 支持文件上传
- 文件列表查看
- 文件删除
- 跨域支持
- 文件大小限制（8MB）

## API接口说明

### 1. 获取文件列表

**请求方法：** GET

**URL：** `/files`

**请求参数：** 无

**返回格式：**
```json
{
    "files": [
        {
            "name": "文件名",
            "url": "/static/文件名"
        }
    ]
}
```

### 2. 上传文件

**请求方法：** POST

**URL：** `/upload`

**请求参数：**
- `file`：文件对象（multipart/form-data）

**返回格式：**
```json
{
    "message": "文件上传成功",
    "url": "/static/文件名"
}
```

### 3. 删除文件

**请求方法：** DELETE

**URL：** `/files/:filename`

**请求参数：**
- `filename`：文件名（URL路径参数）

**返回格式：**
```json
{
    "message": "文件删除成功"
}
```

## 使用示例

1. 启动服务器：
```bash
go run main.go
```

2. 访问地址：
```
http://localhost:8080
```

服务器将在8080端口启动，可以通过浏览器或API工具进行接口测试。
