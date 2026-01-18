import { create } from "zustand";
import { persist } from "zustand/middleware";
import { AuthStore, TDomain, TMenuWidget, TRole, TUser } from "./types";

export const LocalStorageKey_AuthStore = "user_auth";
export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get, store) => ({
      token: null,
      user: null,
      domain: null,
      role: null,
      menus: null,
      widgets: null,
      setAuth: (
        token: string | null,
        user: TUser | null,
        domain: TDomain | null,
        role: TRole | null,
        menus: TMenuWidget[] | null,
        widgets: TMenuWidget[] | null
      ) => set({ token, user, domain, role, menus, widgets }),
      clearAuth: () => {
        useAuthStore.persist.clearStorage();
        set(store.getInitialState());
      },
    }),
    {
      name: LocalStorageKey_AuthStore, // storage key
      // 默认使用 localStorage
    }
  )
);
