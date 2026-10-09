export type Page = "auth" | "rider" | "driver" | "ride" | "notifications";

export const pages: Page[] = ["auth", "rider", "driver", "ride", "notifications"];

export function pageFromPath(pathname: string): Page {
  const segment = pathname.replace(/^\/+/, "").split("/")[0];
  return pages.includes(segment as Page) ? (segment as Page) : "auth";
}

export function pathForPage(page: Page) {
  return `/${page}`;
}
