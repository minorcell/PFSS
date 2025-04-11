# bucker_module

Bucket(桶)：用户自定义的文件夹路径

bucker_module 的主要功能是管理桶的创建、删除、修改等操作。

### 桶

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

### 桶列表

```typescript
interface BucketList {
  buckets: Bucket[];
}
```

