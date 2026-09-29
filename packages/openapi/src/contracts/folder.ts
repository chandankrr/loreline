import { initContract } from "@ts-rest/core";
import { z } from "zod";

import {
  ZFolder,
  ZFolderSummary,
  ZFolderWithCounts,
  ZTreeRecord,
} from "@loreline/zod";

import { getSecurityMetadata } from "../utils.js";

const c = initContract();

const metadata = getSecurityMetadata();

export const folderContract = c.router(
  {
    createFolder: {
      summary: "Create Folder",
      path: "",
      method: "POST",
      description: "Create a new folder",
      body: ZFolder.pick({
        name: true,
        parentId: true,
      })
        .partial()
        .required({
          name: true,
        }),
      responses: {
        201: ZFolder,
      },
      metadata: metadata,
    },

    listFolders: {
      summary: "List Folders",
      path: "",
      method: "GET",
      description: "List folders for a user",
      query: z.object({
        parentId: z.string().uuid().optional(),
      }),
      responses: {
        200: z.object({
          folders: z.array(ZFolderWithCounts),
        }),
      },
      metadata: metadata,
    },

    getFolder: {
      summary: "Get Folder",
      path: "/:id",
      method: "GET",
      description: "Get a folder by its ID",
      responses: {
        200: ZFolderSummary,
      },
      metadata: metadata,
    },

    getFolderBreadcrumb: {
      summary: "Get Folder Breadcrumb",
      path: "/:id/breadcrumb",
      method: "GET",
      description: "Get the breadcrumb path for a folder",
      responses: {
        200: z.object({
          breadcrumb: z.array(ZFolder),
        }),
      },
    },

    getFolderTree: {
      summary: "Get Folder Tree",
      path: "/tree",
      method: "GET",
      description: "Get the folder hierarchy as a tree",
      responses: {
        200: z.object({
          tree: z.array(ZTreeRecord),
        }),
      },
      metadata: metadata,
    },

    updateFolder: {
      summary: "Update Folder",
      path: "/:id",
      method: "PATCH",
      description: "Update a folder",
      body: ZFolder.pick({
        name: true,
      }).partial(),
      responses: {
        200: ZFolder,
      },
      metadata: metadata,
    },

    deleteFolder: {
      summary: "Delete Folder",
      path: "/:id",
      method: "DELETE",
      description: "Delete a folder",
      responses: {
        201: z.void(),
      },
      metadata: metadata,
    },

    moveFolder: {
      summary: "Move Folder",
      path: "/:id/move",
      method: "POST",
      description: "Move a folder to a different destination",
      body: z.object({
        parentId: z.string().uuid().nullish(),
      }),
      responses: {
        200: ZFolder,
      },
      metadata: metadata,
    },
  },
  {
    pathPrefix: "/api/v1/folders",
  },
);
