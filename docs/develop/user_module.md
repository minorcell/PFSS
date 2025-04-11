# user_module

在 PFSS 中，用户分为系统用户和可访问性用户，所有在部署服务器上注册的用户都是系统用户，所有通过互联网链接访问 PFSS 内容的用户都是可访问性用户。用户注册功能仅仅是rootUser（跟用户，全局唯一）可以注册其他用户。

- 系统用户：PFSS 的系统用户，用于管理 PFSS 中的文件，包括所有权限内容
- 可访问性用户：PFSS 的可访问性用户，仅有 READ 权限，用于访问 PFSS 中的文件

### 系统用户

```typescript
interface SystemUser {
  id: string;
  username: string;
  password: string;
  createdAt: number;
  isRootUser: boolean;
}
```


