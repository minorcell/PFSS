# file_module

file_module 的主要功能是管理文件的上传、预览、下载、删除等操作。

### 文件上传

```typescript
interface FileUpload {
  file: File;
  customNameType?: "time" | "filename" | "hash";
  customFolder?: string;
}
```

### 文件预览

```typescript
interface FilePreview {
  id: string;
  name: string;
  path: string;
  folder: string;
  size: number;
  type: string;
  lastModified: number;
}
```

### 文件下载

```typescript
interface FileDownload {
  id: string;
}
```

### 文件删除

```typescript
interface FileDelete {
  id: string;
}
```

### 桶管理

```typescript
interface Bucket {
  id: string;
  name: string;
  fileIds: File.id[];
}
```

### 桶创建

```typescript
interface BucketCreate {
  name: string;
  fileIds: File.id[];
}
```

### 桶删除

```typescript
interface BucketDelete {
  id: string;
}
```

### 桶修改

```typescript
interface BucketUpdate {
  id: string;
  name: string;
  fileIds: File.id[];
}
```
