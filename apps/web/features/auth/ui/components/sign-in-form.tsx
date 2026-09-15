"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";

import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRightIcon, LoaderCircleIcon } from "lucide-react";
import { Controller, useForm } from "react-hook-form";
import type { z } from "zod";

import { Button } from "@loreline/ui/components/button";
import { Field, FieldError, FieldLabel } from "@loreline/ui/components/field";
import { Input } from "@loreline/ui/components/input";
import { toast } from "@loreline/ui/components/toast";

import { getApiErrorCode } from "@/api/utils";
import { getSafeRedirect } from "@/lib/utils";

import { useLogin } from "../../api";
import { signInSchema } from "../../schemas";
import { GoogleIcon } from "../icons/google";

export const SignInForm = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { mutate, isPending } = useLogin();

  const redirectTo = getSafeRedirect(searchParams.get("redirect"));

  const form = useForm<z.infer<typeof signInSchema>>({
    resolver: zodResolver(signInSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  function onSubmit(data: z.infer<typeof signInSchema>) {
    mutate(
      { body: data },
      {
        onSuccess: () => {
          toast.add({ title: "Signed in successfully", type: "success" });
          router.push(redirectTo);
        },
        onError: (error) => {
          if (getApiErrorCode(error) === "EMAIL_NOT_VERIFIED") {
            router.push(
              `/verify-email?email=${encodeURIComponent(data.email)}`,
            );
            return;
          }
        },
      },
    );
  }

  return (
    <div className="m-auto w-full max-w-md py-14">
      <p className="font-semibold text-brand-ink text-sm">Welcome back</p>
      <h2 className="mt-3 font-semibold text-5xl leading-[1.08] tracking-[-0.047em]">
        Return to your books.
      </h2>
      <p className="mt-3 text-muted-foreground text-sm leading-relaxed">
        Sign in to pick up exactly where you left off.
      </p>

      <div className="mt-9 space-y-5">
        <Button type="button" variant="outline" size="xl" className="w-full">
          <GoogleIcon className="size-3.5 grayscale-50" /> Continue with Google
        </Button>
        <div className="flex items-center gap-3 text-muted-foreground text-xs">
          <span className="h-px flex-1 bg-border" />
          <span>or</span>
          <span className="h-px flex-1 bg-border" />
        </div>

        <form
          id="form-sign-in"
          onSubmit={form.handleSubmit(onSubmit)}
          className="space-y-5"
        >
          <Controller
            name="email"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor="form-sign-in-email">Email</FieldLabel>
                <Input
                  {...field}
                  id="form-sign-in-email"
                  aria-invalid={fieldState.invalid}
                  autoComplete="email"
                  placeholder="you@example.com"
                />
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />
          <Controller
            name="password"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <div className="flex items-center justify-between">
                  <FieldLabel htmlFor="form-sign-in-password">
                    Password
                  </FieldLabel>
                  <Link
                    href="/forgot-password"
                    className="text-muted-foreground text-xs underline-offset-1 hover:underline"
                  >
                    Forgot password?
                  </Link>
                </div>
                <Input
                  {...field}
                  id="form-sign-in-password"
                  aria-invalid={fieldState.invalid}
                  type="password"
                  autoComplete="current-password"
                  placeholder="Your password"
                />
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />
          <Button
            type="submit"
            size="xl"
            className="w-full"
            disabled={isPending}
          >
            {isPending ? <LoaderCircleIcon className="animate-spin" /> : null}
            Sign in
            <ArrowRightIcon data-icon="inline-end" />
          </Button>
        </form>
      </div>

      <p className="mt-6 text-center text-muted-foreground text-sm">
        New to Loreline?{" "}
        <Link
          href="/sign-up"
          className="font-semibold text-foreground underline-offset-4 hover:underline"
        >
          Create an account
        </Link>
      </p>
      <p className="mt-10 text-center text-muted-foreground text-xs">
        By continuing, you agree to Loreline&apos;s{" "}
        <Link href="/privacy" className="underline underline-offset-1">
          privacy principles
        </Link>
        .
      </p>
    </div>
  );
};
