export interface Record {
  Student_id: string;
  Student_name: string;
  Subject: string;
  Grade: string;
}

export const SearchParamsKeys = {
  PAGE: "page",
  PAGE_SIZE: "pageSize",
  SORT_BY: "sortBy",
  SORT_ORDER: "sortOrder",
  NAME: "name",
  SUBJECT: "subject",
} as const;
