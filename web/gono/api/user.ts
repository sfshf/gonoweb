import { postJson } from ".";
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
          "utf-8"
        )
      )
    )
  ).toString("hex");
  // 发送请求
  const url = process.env.NEXT_PUBLIC_API_BASE + "/user/signIn";
  const resp: any = await postJson(url, {
    ...req,
    password,
  });
  return resp;
};

export const signOut = async (): Promise<Object> => {
  // 发送请求
  const url = process.env.NEXT_PUBLIC_API_BASE + "/user/signOut";
  const resp: any = await postJson(url);
  return resp;
};
