export interface INavItem {
  href: string;
  name: string;
}

export const navigation: INavItem[] = [
  {
    name: '许可',
    href: '/licences',
  },
  {
    name: '成员',
    href: '/members',
  },
  {
    name: '邀请',
    href: '/invitations'
  },
  {
    name: '交易历史',
    href: '/orders',
  },
  {
    name: '设置',
    href: '/settings',
  },
  {
    name: '退出',
    href: '/logout',
  },
];
