import { deleteJson, getJson, postJson, putJson } from ".";
import { md5 } from "@noble/hashes/legacy.js";

export const visit = async () => {};

export type SignInReq = {
  account: string;
  password: string;
};

export const signIn = async (req: SignInReq): Promise<Object> => {
  if (req.account == "") {
    throw new Error(`account is empty`, {
      cause: "Bad Request",
    });
  }
  if (req.password == "") {
    throw new Error(`password is empty`, {
      cause: "Bad Request",
    });
  }
  // 处理用户密码
  const password = Buffer.from(
    md5(
      Uint8Array.from(
        Buffer.from(
          req.password + process.env.NEXT_PUBLIC_PASSWORD_SALT,
          "utf-8",
        ),
      ),
    ),
  ).toString("hex");
  // 发送请求
  return await postJson(process.env.NEXT_PUBLIC_API_BASE + "/user/signIn", {
    ...req,
    password,
  });
};

export const signOut = async (): Promise<Object> => {
  // 发送请求
  return await postJson(process.env.NEXT_PUBLIC_API_BASE + "/user/signOut");
};

export type ListUserReq = {
  page: number;
  pageSize: number;
  email: string;
  nickname: string;
  realname: string;
};

export const listUser = async (req: ListUserReq): Promise<any> => {
  const url =
    process.env.NEXT_PUBLIC_API_BASE +
    "/user?page=" +
    req.page +
    "&pageSize=" +
    req.pageSize +
    "&email=" +
    req.email +
    "&nickname=" +
    req.nickname +
    "&realname=" +
    req.realname;
  return await getJson(url);
};

export type AddUserReq = {
  email: string;
  nickname: string;
};

export const addUser = async (req: AddUserReq): Promise<any> => {
  return await postJson(process.env.NEXT_PUBLIC_API_BASE + "/user", {
    ...req,
  });
};

export type EditUserReq = {
  xid: string;
  email: string;
  nickname: string;
};

export const editUser = async (req: EditUserReq): Promise<any> => {
  return await putJson(process.env.NEXT_PUBLIC_API_BASE + "/user/" + req.xid, {
    email: req.email,
    nickname: req.nickname,
  });
};

export const deleteUser = async (xid: string): Promise<any> => {
  return await deleteJson(process.env.NEXT_PUBLIC_API_BASE + "/user/" + xid);
};
