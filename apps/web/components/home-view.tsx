"use client";

import Image from "next/image";
import Link from "next/link";
import { useEffect, useRef, useState } from "react";

import {
  ArrowRightIcon,
  AudioLinesIcon,
  BookOpenIcon,
  BrainCircuitIcon,
  EyeIcon,
  ImageIcon,
  LibraryIcon,
  Mic2Icon,
  MousePointer2Icon,
  PauseIcon,
  SparklesIcon,
} from "lucide-react";
import { motion, useInView, useReducedMotion } from "motion/react";

import { Button } from "@loreline/ui/components/button";

import { AppHeader } from "./app-header";
import { Logo } from "./logo";

const steps = [
  {
    number: "01",
    icon: MousePointer2Icon,
    title: "Point",
    copy: "Hover a word, line, diagram, or paragraph. Loreline knows exactly what you mean.",
  },
  {
    number: "02",
    icon: Mic2Icon,
    title: "Ask",
    copy: "Speak naturally. Interrupt, follow up, change direction, or ask for the simpler version.",
  },
  {
    number: "03",
    icon: ImageIcon,
    title: "See",
    copy: "Get an answer in voice while the sideboard draws the idea, scene, or relationship beside you.",
  },
];

export const HomeView = () => {
  return (
    <main className="min-h-screen overflow-hidden bg-background text-foreground">
      <AppHeader variant="marketing" />

      <section className="relative px-5 pt-20 pb-24 sm:px-8 sm:pt-28 sm:pb-32">
        <div className="pointer-events-none absolute inset-x-0 top-0 z-0 h-168 bg-[radial-gradient(circle_at_50%_16%,var(--brand-soft),transparent_58%)]" />
        <div className="relative mx-auto max-w-304">
          <motion.div
            initial={{ opacity: 0, y: 16 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.72, ease: [0.16, 1, 0.3, 1] }}
            className="mx-auto max-w-232 text-center"
          >
            <p className="font-semibold text-brand-ink text-sm">
              A reading companion that sees your page
            </p>
            <h1 className="mt-6 text-balance font-normal text-[3.4rem] leading-[0.98] tracking-[-0.055em] sm:text-[5.8rem]">
              Point to a line.
              <br />
              Ask out loud.
            </h1>
            <p className="mx-auto mt-7 max-w-2xl text-pretty text-base text-muted-foreground leading-7 sm:text-lg sm:leading-8">
              Loreline follows your place on the page, answers in real-time
              voice, and makes difficult ideas visible without turning reading
              into another chat window.
            </p>
            <div className="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
              <Button
                size="xl"
                nativeButton={false}
                render={<Link href="/library" prefetch />}
              >
                Open your first book <ArrowRightIcon data-icon="inline-end" />
              </Button>
              <Button
                size="xl"
                variant="outline"
                nativeButton={false}
                render={<a href="#voice" />}
              >
                <AudioLinesIcon /> Hear how it works
              </Button>
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 26, scale: 0.985 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            transition={{
              duration: 0.9,
              delay: 0.12,
              ease: [0.16, 1, 0.3, 1],
            }}
            className="mt-16 sm:mt-20"
          >
            <HeroReadingRoom />
          </motion.div>
        </div>
      </section>

      <section id="how" className="border-y px-5 py-24 sm:px-8 sm:py-32">
        <div className="mx-auto max-w-304">
          <Reveal className="grid gap-7 lg:grid-cols-[0.8fr_1.2fr] lg:items-end">
            <p className="font-semibold text-brand-ink text-sm">
              The page is the prompt
            </p>
            <h2 className="text-balance font-normal text-4xl leading-[1.04] tracking-[-0.045em] sm:text-6xl">
              No copying passages.
              <br />
              No describing where you are.
            </h2>
          </Reveal>
          <div className="mt-16 grid border-y md:grid-cols-3">
            {steps.map((step, index) => (
              <Reveal
                key={step.title}
                delay={index * 0.07}
                className="border-b p-7 last:border-b-0 sm:p-9 md:border-r md:border-b-0 md:last:border-r-0"
              >
                <div className="flex items-center justify-between">
                  <span className="font-mono text-muted-foreground text-xs">
                    {step.number}
                  </span>
                  <span className="grid size-11 place-items-center rounded-full bg-brand-soft text-brand-ink">
                    <step.icon className="size-5" />
                  </span>
                </div>
                <h3 className="mt-20 font-medium text-3xl tracking-[-0.035em]">
                  {step.title}
                </h3>
                <p className="mt-4 max-w-sm text-muted-foreground leading-7">
                  {step.copy}
                </p>
              </Reveal>
            ))}
          </div>
        </div>
      </section>

      <section
        id="voice"
        className="bg-foreground px-5 py-24 text-background sm:px-8 sm:py-32"
      >
        <div className="mx-auto grid max-w-304 gap-14 lg:grid-cols-[0.8fr_1.2fr] lg:items-center">
          <Reveal>
            <p className="font-semibold text-brand text-sm">Realtime voice</p>
            <h2 className="mt-6 text-pretty font-normal text-4xl leading-[1.03] tracking-[-0.045em] sm:text-6xl">
              Stay with the thought, not the interface.
            </h2>
            <p className="mt-6 max-w-lg text-pretty text-background/62 leading-7">
              Ask mid-sentence. Interrupt the answer. Jump back two pages.
              Loreline keeps the conversation attached to what you can see.
            </p>
          </Reveal>

          <Reveal delay={0.1}>
            <VoiceConversation />
          </Reveal>
        </div>
      </section>

      <section id="sideboard" className="px-5 py-24 sm:px-8 sm:py-36">
        <div className="mx-auto max-w-304">
          <Reveal className="max-w-3xl">
            <p className="font-semibold text-brand-ink text-sm">
              Visual sideboard
            </p>
            <h2 className="mt-6 text-balance font-normal text-4xl leading-[1.03] tracking-[-0.045em] sm:text-6xl">
              Some ideas need to be seen.
            </h2>
            <p className="mt-6 max-w-xl text-muted-foreground leading-7">
              Ask for the mechanism, the metaphor, or the scene. Loreline can
              keep talking while it builds the visual beside your book.
            </p>
          </Reveal>

          <div className="mt-14 grid gap-5 lg:grid-cols-2">
            <Reveal>
              <article className="overflow-hidden rounded-[2rem] bg-card p-3">
                <Image
                  src="/images/loreline-alice-v2.webp"
                  alt="An original interpretation of Alice falling through the rabbit hole"
                  width={1536}
                  height={1024}
                  sizes="(max-width: 1023px) calc(100vw - 40px), 600px"
                  className="aspect-3/2 w-full rounded-[1.35rem] object-cover"
                />
                <div className="p-5 sm:p-6">
                  <p className="font-semibold text-brand-ink text-xs">
                    Alice's Adventures in Wonderland
                  </p>
                  <h3 className="mt-3 font-medium text-2xl tracking-[-0.03em]">
                    Fall through the passage, not just past it.
                  </h3>
                </div>
              </article>
            </Reveal>

            <Reveal delay={0.08}>
              <article className="overflow-hidden rounded-[2rem] bg-foreground p-3 text-background">
                <Image
                  src="/images/loreline-gatsby-v2.webp"
                  alt="An original interpretation of Gatsby looking across the bay toward the green light"
                  width={1536}
                  height={1024}
                  sizes="(max-width: 1023px) calc(100vw - 40px), 600px"
                  className="aspect-3/2 w-full rounded-[1.35rem] object-cover"
                />
                <div className="p-5 sm:p-6">
                  <p className="font-semibold text-brand text-xs">
                    The Great Gatsby
                  </p>
                  <h3 className="mt-3 font-medium text-2xl tracking-[-0.03em]">
                    See how much longing fits inside one green light.
                  </h3>
                </div>
              </article>
            </Reveal>
          </div>
        </div>
      </section>

      <section className="border-y bg-card px-5 py-24 sm:px-8 sm:py-32">
        <div className="mx-auto max-w-304">
          <Reveal className="mx-auto mb-14 max-w-3xl text-center">
            <p className="font-semibold text-brand-ink text-sm">
              Reading that stays
            </p>
            <h2 className="mt-6 text-balance font-normal text-4xl leading-[1.03] tracking-[-0.045em] sm:text-6xl">
              A page turns. The idea does not disappear.
            </h2>
            <p className="mx-auto mt-6 max-w-xl text-muted-foreground leading-7">
              Loreline keeps the scene, question, and explanation attached to
              the book—so understanding accumulates instead of vanishing into
              chat history.
            </p>
          </Reveal>
          <Reveal delay={0.08}>
            <PageToMind />
          </Reveal>
        </div>
      </section>

      <section className="bg-sky-soft px-5 py-24 sm:px-8 sm:py-32">
        <div className="mx-auto max-w-304">
          <Reveal className="grid gap-12 lg:grid-cols-[0.9fr_1.1fr] lg:items-center">
            <div>
              <p className="font-semibold text-sky text-sm">
                Context, in order
              </p>
              <h2 className="mt-6 text-balance font-normal text-4xl leading-[1.03] tracking-[-0.045em] sm:text-6xl">
                Your page first. Retrieval second.
              </h2>
              <p className="mt-6 max-w-lg text-muted-foreground leading-7">
                Loreline begins with the passage under your pointer. It reaches
                into the chapter or the rest of the book only when the answer
                actually needs it.
              </p>
            </div>
            <div className="rounded-[2rem] bg-background p-4 shadow-soft sm:p-6">
              {[
                [
                  "Visible page",
                  "Primary context",
                  "bg-brand text-primary-foreground",
                ],
                ["Current chapter", "When the thought reaches back", "bg-card"],
                [
                  "Whole book",
                  "Only when broader retrieval helps",
                  "bg-control",
                ],
              ].map(([title, copy, color], index) => (
                <div
                  key={title}
                  className={`flex items-center gap-4 rounded-[1.2rem] p-5 ${color} ${index ? "mt-3" : ""}`}
                >
                  <span className="font-mono text-xs opacity-65">
                    0{index + 1}
                  </span>
                  <div>
                    <p className="font-semibold">{title}</p>
                    <p className="mt-1 text-xs opacity-65">{copy}</p>
                  </div>
                </div>
              ))}
            </div>
          </Reveal>
        </div>
      </section>

      <section className="px-5 py-20 sm:px-8 sm:py-28">
        <Reveal className="mx-auto max-w-304 overflow-hidden rounded-[2.25rem] bg-brand p-8 text-primary-foreground sm:p-14 lg:p-20">
          <div className="grid gap-12 lg:grid-cols-[1fr_auto] lg:items-end">
            <div>
              <span className="grid size-12 place-items-center rounded-full bg-primary-foreground text-brand">
                <LibraryIcon className="size-5" />
              </span>
              <h2 className="mt-10 max-w-4xl text-balance font-normal text-5xl leading-[0.98] tracking-tighter sm:text-7xl">
                Bring a book.
                <br />
                Leave with a world.
              </h2>
              <p className="mt-6 max-w-xl text-primary-foreground/72 leading-7">
                Your library, reading progress, conversations, and visual
                explanations stay connected to the books that inspired them.
              </p>
            </div>
            <Button
              size="xl"
              variant="secondary"
              className="bg-background text-foreground hover:bg-background/90"
              nativeButton={false}
              render={<Link href="/library" />}
            >
              Start with your first book{" "}
              <ArrowRightIcon data-icon="inline-end" />
            </Button>
          </div>
        </Reveal>
      </section>

      <footer className="border-t px-5 py-12 sm:px-8">
        <div className="mx-auto flex max-w-304 flex-col gap-8 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex flex-col items-center gap-2 sm:items-start">
            <Logo />
            <p className="mt-3 text-muted-foreground text-sm">
              Stay in the book. Go deeper.
            </p>
          </div>
          <div className="flex flex-wrap justify-center gap-x-6 gap-y-3 text-muted-foreground text-sm sm:justify-end">
            <Link href="/library" className="hover:text-foreground">
              Library
            </Link>
            <Link href="/privacy" className="hover:text-foreground">
              Privacy
            </Link>
            <a
              href="mailto:hello@loreline.app"
              className="hover:text-foreground"
            >
              Contact
            </a>
          </div>
        </div>
      </footer>
    </main>
  );
};

const Reveal = ({
  children,
  className,
  delay = 0,
}: {
  children: React.ReactNode;
  className?: string;
  delay?: number;
}) => {
  const reduce = useReducedMotion();

  return (
    <motion.div
      initial={{ opacity: 0, y: reduce ? 0 : 22 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, amount: 0.16 }}
      transition={{
        duration: reduce ? 0 : 0.72,
        delay: reduce ? 0 : delay,
        ease: [0.16, 1, 0.3, 1],
      }}
      className={className}
    >
      {children}
    </motion.div>
  );
};

export const HeroReadingRoom = () => {
  return (
    <div
      data-layout-card="hero-reader"
      className="overflow-hidden rounded-[1.75rem] border bg-background shadow-float sm:rounded-[2.25rem]"
    >
      <div className="flex h-14 items-center justify-between border-b px-4 sm:px-6">
        <div className="flex items-center gap-3">
          <span className="grid size-8 place-items-center rounded-full bg-foreground text-background">
            <BookOpenIcon className="size-4" />
          </span>
          <div>
            <p className="font-semibold text-xs">Moby-Dick</p>
            <p className="text-[0.65rem] text-muted-foreground">
              Page 507 of 635
            </p>
          </div>
        </div>
        <span className="hidden items-center gap-2 rounded-full bg-sky-soft px-3 py-1.5 font-medium text-sky text-xs sm:inline-flex">
          <span className="size-1.5 rounded-full bg-sky" /> Seeing this page
        </span>
      </div>

      <div className="grid lg:grid-cols-[minmax(0,1fr)_23rem]">
        <div className="relative flex min-h-144 items-center justify-center bg-canvas p-6 sm:p-10">
          <ReadingPage />
          <div className="absolute bottom-5 left-1/2 flex w-[calc(100%-2.5rem)] max-w-md -translate-x-1/2 items-center gap-3 rounded-full bg-foreground p-2.5 pl-4 text-background shadow-panel sm:bottom-7">
            <span className="grid size-9 shrink-0 place-items-center rounded-full bg-brand text-primary-foreground">
              <Mic2Icon className="size-4" />
            </span>
            <div className="min-w-0 flex-1">
              <p className="truncate font-medium text-xs">
                “Put me on the deck when the whale appears.”
              </p>
              <p className="text-[0.65rem] text-background/60">
                Listening · interrupt anytime
              </p>
            </div>
            <div className="hidden text-brand sm:block">
              <VoiceBars compact />
            </div>
          </div>
        </div>

        <aside className="flex min-h-144 flex-col border-t bg-background lg:border-t-0 lg:border-l">
          <div className="border-b p-5">
            <div className="flex items-center justify-between">
              <div>
                <p className="font-semibold text-sm">Visual sideboard</p>
                <p className="mt-1 text-muted-foreground text-xs">
                  Built from the passage
                </p>
              </div>
              <SparklesIcon className="size-4 text-brand-ink" />
            </div>
          </div>
          <div className="flex flex-1 flex-col gap-4 p-4">
            <p className="rounded-[1.1rem] rounded-bl-sm bg-card p-4 text-muted-foreground text-sm leading-6">
              The whale arrives as animal, obsession, and fate at once—so vast
              that it seems to consume the frame around Ahab.
            </p>
            <div className="overflow-hidden rounded-[1.2rem] border bg-card">
              <Image
                src="/images/loreline-moby-dick-v2.webp"
                alt="An original interpretation of the white whale surfacing beside the Pequod"
                width={1536}
                height={1024}
                sizes="(max-width: 1023px) calc(100vw - 72px), 336px"
                className="aspect-4/3 w-full object-cover"
                priority
              />
              <div className="flex items-center gap-2 border-t bg-background p-3 text-muted-foreground text-xs">
                <EyeIcon className="size-3.5 text-brand-ink" />
                The whale breaks the surface
              </div>
            </div>
            <div className="mt-auto flex gap-2">
              <span className="rounded-full border px-3 py-1.5 text-xs">
                Simplify
              </span>
              <span className="rounded-full border px-3 py-1.5 text-xs">
                Another visual
              </span>
            </div>
          </div>
        </aside>
      </div>
    </div>
  );
};

const ReadingPage = ({ compact = false }: { compact?: boolean }) => {
  return (
    <div
      className={`reader-paper relative mx-auto w-full max-w-104 overflow-hidden rounded-[1.1rem] border border-reader-ink/10 bg-reader-paper text-reader-ink ${
        compact ? "min-h-100 p-7" : "min-h-136 p-8 sm:p-11"
      }`}
    >
      <div className="flex items-center justify-between font-semibold text-[0.62rem] text-reader-muted tracking-[0.12em]">
        <span>Chapter 135</span>
        <span>507</span>
      </div>
      <p className="mt-12 text-center text-reader-muted text-xs tracking-[0.12em]">
        The Chase — Third Day
      </p>
      <div className="mx-auto mt-9 max-w-[18rem] space-y-5 text-center text-reader-muted text-sm leading-7">
        <p>The sea rolled on as it had rolled five thousand years ago.</p>
        <p className="relative rounded-lg bg-reader-highlight px-2 py-1.5 text-reader-ink">
          Towards thee I roll, thou all-destroying but unconquering whale.
          <span className="absolute -top-2 -right-2 grid size-6 place-items-center rounded-full bg-brand text-primary-foreground shadow-soft">
            <MousePointer2Icon className="size-3.5 fill-current" />
          </span>
        </p>
        <p>
          The white body rose through the dark water until the horizon seemed
          too small to hold it.
        </p>
      </div>
    </div>
  );
};

const VoiceBars = ({ compact = false }: { compact?: boolean }) => {
  const reduce = useReducedMotion();
  const heights = compact
    ? [8, 13, 19, 25, 31, 25, 19, 13, 8]
    : [12, 20, 32, 46, 60, 72, 60, 46, 32, 20, 12];

  return (
    <div
      className="flex items-center justify-center gap-1.5"
      aria-label="Live voice activity"
    >
      {heights.map((height, index) => (
        <motion.span
          // biome-ignore lint/suspicious/noArrayIndexKey: Using index as key is acceptable here since the heights array is static and will not change between renders.
          key={`${height}-${index}`}
          className="w-1.5 rounded-full bg-current"
          style={{ height }}
          animate={
            reduce
              ? undefined
              : {
                  scaleY: [0.42, 1, 0.58, 0.84, 0.42],
                  opacity: [0.62, 1, 0.72, 0.9, 0.62],
                }
          }
          transition={{
            duration: 1.8,
            repeat: Number.POSITIVE_INFINITY,
            delay: Math.abs(index - Math.floor(heights.length / 2)) * 0.07,
            ease: "easeInOut",
          }}
        />
      ))}
    </div>
  );
};

const LiveTranscript = ({
  active,
  delay,
  text,
}: {
  active: boolean;
  delay: number;
  text: string;
}) => {
  const reduce = useReducedMotion();
  const [visibleText, setVisibleText] = useState("");

  useEffect(() => {
    if (!active || reduce) return;

    let intervalId: number | undefined;
    const timeoutId = window.setTimeout(() => {
      let index = 0;
      intervalId = window.setInterval(() => {
        index += 1;
        setVisibleText(text.slice(0, index));
        if (index >= text.length && intervalId) {
          window.clearInterval(intervalId);
        }
      }, 23);
    }, delay);

    return () => {
      window.clearTimeout(timeoutId);
      if (intervalId) window.clearInterval(intervalId);
    };
  }, [active, delay, reduce, text]);

  const displayedText = reduce && active ? text : visibleText;
  const typing = active && !reduce && displayedText.length < text.length;

  return (
    <span aria-label={text}>
      <span aria-hidden="true">
        {displayedText}
        {typing && (
          <motion.span
            animate={{ opacity: [1, 0, 1] }}
            transition={{ duration: 0.8, repeat: Number.POSITIVE_INFINITY }}
            className="ml-0.5 inline-block h-[1em] w-px translate-y-[0.12em] bg-current"
          />
        )}
      </span>
    </span>
  );
};

export const VoiceConversation = () => {
  const conversationRef = useRef<HTMLDivElement>(null);
  const active = useInView(conversationRef, { once: true, amount: 0.35 });

  return (
    <div
      ref={conversationRef}
      className="rounded-[2rem] bg-background p-5 text-foreground sm:p-8"
    >
      <div className="flex items-center justify-between border-b pb-5">
        <div className="flex items-center gap-3">
          <span className="relative grid size-11 place-items-center rounded-full bg-brand text-primary-foreground">
            <Mic2Icon className="size-5" />
            <span className="absolute -top-0.5 -right-0.5 size-3 rounded-full border-2 border-background bg-sky" />
          </span>
          <div>
            <p className="font-semibold text-sm">Loreline is with you</p>
            <p className="text-muted-foreground text-xs">
              Live voice · The Creative Act, page 42
            </p>
          </div>
        </div>
        <Button variant="secondary" size="icon" aria-label="Pause voice">
          <PauseIcon />
        </Button>
      </div>

      <div className="min-h-16 py-5 text-muted-foreground text-sm leading-6">
        <LiveTranscript
          active={active}
          delay={250}
          text="Welcome to Loreline. I'm on page 42 with you—what are you curious about?"
        />
      </div>

      <div className="py-8 text-brand sm:py-10">
        <VoiceBars />
      </div>

      <div className="space-y-3" aria-live="polite">
        <div className="ml-auto min-h-18 max-w-[88%] rounded-[1.25rem] rounded-br-sm bg-foreground p-4 text-background text-sm leading-6">
          <LiveTranscript
            active={active}
            delay={2300}
            text="Wait—why does the author call attention a current?"
          />
        </div>
        <div className="min-h-22 max-w-[92%] rounded-[1.25rem] rounded-bl-sm bg-card p-4 text-muted-foreground text-sm leading-6">
          <LiveTranscript
            active={active}
            delay={4000}
            text="Because a current already has direction. You do not create it by force; you notice where it is carrying the work."
          />
        </div>
      </div>

      <div className="mt-5 flex items-center justify-center gap-2 text-muted-foreground text-xs">
        <span className="size-1.5 rounded-full bg-sky" />
        Speak naturally · interrupt anytime
      </div>
    </div>
  );
};

const TurningBook = () => {
  const reduce = useReducedMotion();
  const lines = ["w-4/5", "w-full", "w-3/4", "w-11/12", "w-2/3"];

  return (
    <div className="mx-auto w-full max-w-88">
      <div className="perspective-distant relative aspect-[1.45/1]">
        <div className="absolute inset-x-2 top-2 bottom-0 rounded-[1rem] bg-foreground/12" />
        <div className="absolute top-0 bottom-2 left-0 w-1/2 rounded-l-[1rem] border bg-reader-paper p-5 shadow-soft">
          <p className="font-medium text-reader-ink text-xs">Moby-Dick</p>
          <div className="mt-7 space-y-3">
            {lines.map((width, index) => (
              <span
                // biome-ignore lint/suspicious/noArrayIndexKey: Using index as key is acceptable here since the lines array is static and will not change between renders.
                key={`left-${index}`}
                className={`block h-1 rounded-full bg-reader-ink/12 ${width}`}
              />
            ))}
          </div>
        </div>
        <div className="absolute top-0 right-0 bottom-2 w-1/2 rounded-r-[1rem] border bg-reader-paper p-5">
          <div className="space-y-3 pt-8">
            {lines.reverse().map((width, index) => (
              <span
                // biome-ignore lint/suspicious/noArrayIndexKey: Using index as key is acceptable here since the lines array is static and will not change between renders.
                key={`under-${index}`}
                className={`block h-1 rounded-full bg-reader-ink/10 ${width}`}
              />
            ))}
          </div>
        </div>
        <motion.div
          className="backface-hidden absolute top-0 right-0 bottom-2 left-1/2 origin-left rounded-r-[1rem] border bg-reader-paper p-5"
          animate={
            reduce
              ? undefined
              : { rotateY: [0, -8, -176, -176, 0], skewY: [0, -1, 0, 0, 0] }
          }
          transition={{
            duration: 6.4,
            repeat: Number.POSITIVE_INFINITY,
            times: [0, 0.12, 0.44, 0.78, 1],
            ease: [0.45, 0, 0.2, 1],
          }}
        >
          <p className="text-right font-mono text-[0.6rem] text-reader-muted">
            507
          </p>
          <div className="mt-6 space-y-3">
            {lines.map((width, index) => (
              <span
                // biome-ignore lint/suspicious/noArrayIndexKey: Using index as key is acceptable here since the lines array is static and will not change between renders.
                key={`turn-${index}`}
                className={`ml-auto block h-1 rounded-full bg-reader-ink/14 ${width}`}
              />
            ))}
            <span className="ml-auto block h-5 w-4/5 rounded bg-reader-highlight" />
          </div>
        </motion.div>
      </div>
      <p className="mt-5 text-center text-muted-foreground text-sm">
        The page moves on
      </p>
    </div>
  );
};

const KnowledgeTrack = () => {
  const reduce = useReducedMotion();
  const fragments = ["the scene", "the question", "the connection"];

  return (
    <>
      <div className="relative hidden h-52 md:block">
        <span className="absolute top-1/2 right-4 left-4 h-px bg-border" />
        {fragments.map((fragment, index) => (
          <motion.span
            key={fragment}
            className="absolute top-1/2 left-3 whitespace-nowrap rounded-full border bg-background px-3 py-2 text-xs shadow-soft"
            animate={
              reduce
                ? { x: 110, y: index * 34 - 34, opacity: 1 }
                : {
                    x: [-20, 80, 170],
                    y: [index * 24 - 24, index * 36 - 52, index * 24 - 24],
                    opacity: [0, 1, 0],
                    scale: [0.92, 1, 0.96],
                  }
            }
            transition={{
              duration: 4.8,
              repeat: Number.POSITIVE_INFINITY,
              delay: index * 0.7,
              ease: [0.45, 0, 0.2, 1],
            }}
          >
            {fragment}
          </motion.span>
        ))}
      </div>

      <div className="relative h-36 md:hidden">
        <span className="absolute top-3 bottom-3 left-1/2 w-px bg-border" />
        {fragments.map((fragment, index) => (
          <motion.span
            key={`mobile-${fragment}`}
            className="absolute top-0 left-1/2 -translate-x-1/2 whitespace-nowrap rounded-full border bg-background px-3 py-2 text-xs shadow-soft"
            animate={
              reduce
                ? { y: 48 + index * 6, opacity: 1 }
                : {
                    x: [index * 8 - 8, index * -7 + 7, 0],
                    y: [-10, 48, 106],
                    opacity: [0, 1, 0],
                    scale: [0.92, 1, 0.96],
                  }
            }
            transition={{
              duration: 4.5,
              repeat: Number.POSITIVE_INFINITY,
              delay: index * 0.65,
              ease: [0.45, 0, 0.2, 1],
            }}
          >
            {fragment}
          </motion.span>
        ))}
      </div>
    </>
  );
};

const PageToMind = () => {
  return (
    <div className="grid items-center gap-8 rounded-[2rem] border bg-background p-6 shadow-soft sm:p-10 md:grid-cols-[1fr_14rem_1fr] lg:p-14">
      <TurningBook />
      <KnowledgeTrack />
      <div className="mx-auto w-full max-w-88 rounded-[1.6rem] bg-foreground p-6 text-background shadow-panel sm:p-8">
        <span className="grid size-12 place-items-center rounded-full bg-brand text-primary-foreground">
          <BrainCircuitIcon className="size-6" />
        </span>
        <p className="mt-8 font-medium text-2xl tracking-[-0.03em]">
          What stays with you
        </p>
        <div className="mt-5 space-y-2 text-sm">
          {[
            "The scene you pictured",
            "The question you followed",
            "The idea in your own words",
          ].map((item, index) => (
            <div
              key={item}
              className="flex items-center gap-3 rounded-[0.9rem] bg-background/8 p-3 text-background/72"
            >
              <span className="font-mono text-[0.62rem] text-brand">
                0{index + 1}
              </span>
              {item}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
