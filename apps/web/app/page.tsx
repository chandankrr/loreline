import { Button } from "@loreline/ui/components/button";

export default function Page() {
  return (
    <div className="flex min-h-svh p-6">
      <div className="flex min-w-0 max-w-md flex-col gap-4 text-sm leading-loose">
        <div>
          <h1 className="font-medium">Loreline</h1>
          <p>You may now add components and start building.</p>
          <Button className="mt-2">Button</Button>
        </div>
        <div className="font-mono text-muted-foreground text-xs">
          (Press <kbd>d</kbd> to toggle dark mode)
        </div>
      </div>
    </div>
  );
}
