import { deleteJson, getJson, postJson, putJson } from ".";

export type ListResourceReq = {
  page: number;
  pageSize: number;
  type: number;
  identifier: string;
  name: string;
};

export const listResource = async (req: ListResourceReq): Promise<any> => {
  const url =
    process.env.NEXT_PUBLIC_API_BASE +
    "/resource?page=" +
    req.page +
    "&pageSize=" +
    req.pageSize +
    "&type=" +
    req.type +
    "&identifier=" +
    req.identifier +
    "&name=" +
    req.name;
  return await getJson(url);
};

export type AddResourceReq = {
  type: number;
  identifier: string;
  name: string;
  intro: string;
  icon: string;
};

export const addResource = async (req: AddResourceReq): Promise<any> => {
  return await postJson(process.env.NEXT_PUBLIC_API_BASE + "/resource", {
    ...req,
  });
};

export type EditResourceReq = {
  id: number;
  name: string;
  type: number;
  identifier: string;
  intro: string;
  icon: string;
};

export const editResource = async (req: EditResourceReq): Promise<any> => {
  return await putJson(
    process.env.NEXT_PUBLIC_API_BASE + "/resource/" + req.id,
    {
      name: req.name,
      intro: req.intro,
      icon: req.icon,
    },
  );
};

export const deleteResource = async (id: number): Promise<any> => {
  return await deleteJson(process.env.NEXT_PUBLIC_API_BASE + "/resource/" + id);
};
