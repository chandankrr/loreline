import { AudioLinesIcon } from "lucide-react";

import { Badge } from "@loreline/ui/components/badge";
import { cn } from "@loreline/ui/lib/utils";

import {
  LORELINE_REALTIME_MODEL_ID,
  LORELINE_REALTIME_MODEL_LABEL,
} from "@/lib/constants";

export const RealtimeModelBadge = ({ className }: { className?: string }) => {
  return (
    <Badge
      variant="outline"
      title={`OpenAI model: ${LORELINE_REALTIME_MODEL_ID}`}
      aria-label={`Voice model: ${LORELINE_REALTIME_MODEL_LABEL}`}
      className={cn(
        "h-6 gap-1.5 border-sky/35 bg-sky-soft/55 px-2 font-semibold text-[0.65rem] text-foreground shadow-sm",
        className,
      )}
    >
      <AudioLinesIcon className="text-sky" />
      {LORELINE_REALTIME_MODEL_LABEL}
    </Badge>
  );
};
