import { deleteJson, getJson, postJson, putJson } from ".";
import { md5 } from "@noble/hashes/legacy.js";

export type ListRoleReq = {
  page: number;
  pageSize: number;
  name: string;
};

export const listRole = async (req: ListRoleReq): Promise<any> => {
  const url =
    process.env.NEXT_PUBLIC_API_BASE +
    "/role?page=" +
    req.page +
    "&pageSize=" +
    req.pageSize +
    "&name=" +
    req.name;
  return await getJson(url);
};

export type AddRoleReq = {
  name: string;
  intro: string;
};

export const addRole = async (req: AddRoleReq): Promise<any> => {
  return await postJson(process.env.NEXT_PUBLIC_API_BASE + "/role", {
    ...req,
  });
};

export type EditRoleReq = {
  xid: string;
  name: string;
  intro: string;
};

export const editRole = async (req: EditRoleReq): Promise<any> => {
  return await putJson(process.env.NEXT_PUBLIC_API_BASE + "/role/" + req.xid, {
    name: req.name,
    intro: req.intro,
  });
};

export const deleteRole = async (xid: string): Promise<any> => {
  return await deleteJson(process.env.NEXT_PUBLIC_API_BASE + "/role/" + xid);
};
