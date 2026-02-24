import { useAuthStore } from "@/zustand/user";

export const postJson = async (url: string, json?: Object): Promise<any> => {
  const token = useAuthStore.getState().token;
  const opt = {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    ...(json ? { body: JSON.stringify(json) } : {}),
  };
  const resp = await fetch(url, opt);
  const data: any = await resp.json();
  if (resp.status >= 500) {
    throw new Error(`${data.msg}`, { cause: "Internal Error" });
  } else if (resp.status >= 400) {
    if (resp.status == 401) {
      useAuthStore.getState().clearAuth();
    }
    throw new Error(`${data.msg}`, { cause: "Bad Request" });
  } else if (resp.ok) {
    return data;
  }
};

export const getJson = async (url: string): Promise<any> => {
  const token = useAuthStore.getState().token;
  const opt = {
    method: "GET",
    headers: {
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  };
  const resp = await fetch(url, opt);
  const data: any = await resp.json();
  if (resp.status >= 500) {
    throw new Error(`${data.msg}`, { cause: "Internal Error" });
  } else if (resp.status >= 400) {
    if (resp.status == 401) {
      useAuthStore.getState().clearAuth();
    }
    throw new Error(`${data.msg}`, { cause: "Bad Request" });
  } else if (resp.ok) {
    return data;
  }
};

export const putJson = async (url: string, json?: Object): Promise<any> => {
  const token = useAuthStore.getState().token;
  const opt = {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    ...(json ? { body: JSON.stringify(json) } : {}),
  };
  const resp = await fetch(url, opt);
  const data: any = await resp.json();
  if (resp.status >= 500) {
    throw new Error(`${data.msg}`, { cause: "Internal Error" });
  } else if (resp.status >= 400) {
    if (resp.status == 401) {
      useAuthStore.getState().clearAuth();
    }
    throw new Error(`${data.msg}`, { cause: "Bad Request" });
  } else if (resp.ok) {
    return data;
  }
};

export const deleteJson = async (url: string): Promise<any> => {
  const token = useAuthStore.getState().token;
  const opt = {
    method: "DELETE",
    headers: {
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  };
  const resp = await fetch(url, opt);
  const data: any = await resp.json();
  if (resp.status >= 500) {
    throw new Error(`${data.msg}`, { cause: "Internal Error" });
  } else if (resp.status >= 400) {
    if (resp.status == 401) {
      useAuthStore.getState().clearAuth();
    }
    throw new Error(`${data.msg}`, { cause: "Bad Request" });
  } else if (resp.ok) {
    return data;
  }
};
