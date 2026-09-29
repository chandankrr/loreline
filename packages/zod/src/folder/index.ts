import { z } from "zod";

export const ZFolder = z.object({
  id: z.string().uuid(),
  name: z.string(),
  userId: z.string().uuid(),
  parentId: z.string().uuid().nullable(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const ZFolderWithCounts = ZFolder.extend({
  bookCount: z.number(),
  childFolderCount: z.number(),
});

export const ZFolderSummary = ZFolder.extend({
  descendantBookCount: z.number(),
  descendantFolderCount: z.number(),
});

type ZTreeRecord = z.infer<typeof ZFolder> & {
  children: ZTreeRecord[];
};

export const ZTreeRecord: z.ZodType<ZTreeRecord> = ZFolder.extend({
  children: z.lazy(() => z.array(ZTreeRecord)),
});
