import { deleteJson, getJson, postJson, putJson } from ".";
import { md5 } from "@noble/hashes/legacy.js";

export type ListDomainReq = {
  page: number;
  pageSize: number;
  name: string;
};

export const listDomain = async (req: ListDomainReq): Promise<any> => {
  const url =
    process.env.NEXT_PUBLIC_API_BASE +
    "/domain?page=" +
    req.page +
    "&pageSize=" +
    req.pageSize +
    "&name=" +
    req.name;
  return await getJson(url);
};

export type AddDomainReq = {
  email: string;
  intro: string;
};

export const addDomain = async (req: AddDomainReq): Promise<any> => {
  return await postJson(process.env.NEXT_PUBLIC_API_BASE + "/domain", {
    ...req,
  });
};

export type EditDomainReq = {
  xid: string;
  name: string;
  intro: string;
};

export const editDomain = async (req: EditDomainReq): Promise<any> => {
  return await putJson(
    process.env.NEXT_PUBLIC_API_BASE + "/domain/" + req.xid,
    {
      name: req.name,
      intro: req.intro,
    },
  );
};

export const deleteDomain = async (xid: string): Promise<any> => {
  return await deleteJson(process.env.NEXT_PUBLIC_API_BASE + "/domain/" + xid);
};
