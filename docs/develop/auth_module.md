# auth_module

auth_module 的作用不仅仅是区别系统用户和普通用户，还包括文件上传时对 secretKey 的管理。secretKey 是用于文件上传的密钥，在设置 secretKey 时，可以选择时效。

### secretKey

```typescript
interface SecretKey {
  id: string;
  key: string;
  createdAt: number;
  expiresAt: number;
}
```