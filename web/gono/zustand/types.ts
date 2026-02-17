export type AuthStore = {
  token: string | null;
  user: TUser | null;
  domain: TDomain | null;
  role: TRole | null;
  menus: TResource[] | null;
  widgets: TResource[] | null;
  setAuth: (
    token: string | null,
    user: TUser | null,
    domain: TDomain | null,
    role: TRole | null,
    menus: TResource[] | null,
    widgets: TResource[] | null,
  ) => void;
  clearAuth: () => void;
};

export type TUser = {
  id: number;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  xid: string;
  email: string;
  nick_name: string;
  real_name: string;
  password: string;
  avatar: string;
};

export type TDomain = {
  id: number;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  xid: string;
  name: string;
  intro: string;
};

export type TRole = {
  id: number;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  xid: string;
  name: string;
  intro: string;
};

export type TResource = {
  id: number;
  created_at: string;
  updated_at: string;
  deleted_at: string | null;
  type: number;
  identifier: string;
  name: string;
  intro: string;
  icon: string;
};
