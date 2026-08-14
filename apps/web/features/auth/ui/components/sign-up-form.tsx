"use client";

import Link from "next/link";

import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRightIcon } from "lucide-react";
import { Controller, useForm } from "react-hook-form";
import type { z } from "zod";

import { Button } from "@loreline/ui/components/button";
import { Field, FieldError, FieldLabel } from "@loreline/ui/components/field";
import { Input } from "@loreline/ui/components/input";
import { toast } from "@loreline/ui/components/toast";

import { signUpSchema } from "../../schemas";
import { GoogleIcon } from "../icons/google";

export const SignUpForm = () => {
  const form = useForm<z.infer<typeof signUpSchema>>({
    resolver: zodResolver(signUpSchema),
    defaultValues: {
      name: "",
      email: "",
      password: "",
    },
  });

  function onSubmit(data: z.infer<typeof signUpSchema>) {
    console.log(data);
    toast.add({
      title: "Signed up successfully",
      type: "success",
    });
  }

  return (
    <div className="m-auto w-full max-w-md py-14">
      <p className="font-semibold text-brand-ink text-sm">Begin your library</p>
      <h2 className="mt-3 font-semibold text-5xl leading-[1.08] tracking-[-0.047em]">
        Make reading feel alive.
      </h2>
      <p className="mt-3 text-muted-foreground text-sm leading-relaxed">
        Create an account. Your first book is a few seconds away.
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
          id="form-sign-up"
          onSubmit={form.handleSubmit(onSubmit)}
          className="space-y-5"
        >
          <Controller
            name="name"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor="form-sign-up-name">Name</FieldLabel>
                <Input
                  {...field}
                  id="form-sign-up-name"
                  aria-invalid={fieldState.invalid}
                  autoComplete="name"
                  placeholder="How should Loreline address you?"
                />
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />
          <Controller
            name="email"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor="form-sign-up-email">Email</FieldLabel>
                <Input
                  {...field}
                  id="form-sign-up-email"
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
                <FieldLabel htmlFor="form-sign-up-password">
                  Password
                </FieldLabel>
                <Input
                  {...field}
                  id="form-sign-up-password"
                  aria-invalid={fieldState.invalid}
                  type="password"
                  autoComplete="new-password"
                  placeholder="At least 6 characters"
                />
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />
          <Button type="submit" size="xl" className="w-full">
            Sign up
            <ArrowRightIcon data-icon="inline-end" />
          </Button>
        </form>
      </div>

      <p className="mt-6 text-center text-muted-foreground text-sm">
        Already have a library?{" "}
        <Link
          href="/sign-in"
          className="font-semibold text-foreground underline-offset-4 hover:underline"
        >
          Sign in
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
