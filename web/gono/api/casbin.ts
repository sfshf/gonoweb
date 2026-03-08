import { deleteJson, getJson, postJson, putJson } from ".";

export const domainRoles = async (dxid: string): Promise<any> => {
  const url =
    process.env.NEXT_PUBLIC_API_BASE + "/casbin/domain/" + dxid + "/role";
  return await getJson(url);
};

export const domainRoleResources = async (
  dxid: string,
  rxid: string,
): Promise<any> => {
  const url =
    process.env.NEXT_PUBLIC_API_BASE +
    "/casbin/domain/" +
    dxid +
    "/role/" +
    rxid +
    "/resource";
  return await getJson(url);
};

export type AllocDomainRoleResourcesReq = {
  dxid: string;
  rxid: string;
  identifiers: string[];
};

export const allocDomainRoleResources = async (
  req: AllocDomainRoleResourcesReq,
): Promise<any> => {
  const url =
    process.env.NEXT_PUBLIC_API_BASE +
    "/casbin/domain/" +
    req.dxid +
    "/role/" +
    req.rxid +
    "/resource";
  return await postJson(url, {
    identifiers: req.identifiers,
  });
};

export type AllocRoleInDomainReq = {
  xid: string; //用户xid
  dxid: string;
  rxids: string[];
};

export const allocRoleInDomain = async (
  req: AllocRoleInDomainReq,
): Promise<any> => {
  const url = process.env.NEXT_PUBLIC_API_BASE + "/casbin/user/" + req.xid;
  return await postJson(url, {
    domain: req.dxid,
    roles: req.rxids,
  });
};

export const userDomains = async (xid: string): Promise<any> => {
  const url =
    process.env.NEXT_PUBLIC_API_BASE + "/casbin/user/" + xid + "/domain";
  return await getJson(url);
};

export const userRolesInDomain = async (
  xid: string,
  dxid: string,
): Promise<any> => {
  const url =
    process.env.NEXT_PUBLIC_API_BASE +
    "/casbin/user/" +
    xid +
    "/domain/" +
    dxid +
    "/role";
  return await getJson(url);
};
