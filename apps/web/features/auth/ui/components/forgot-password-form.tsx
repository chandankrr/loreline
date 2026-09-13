"use client";

import Link from "next/link";

import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRightIcon, LoaderCircleIcon } from "lucide-react";
import { Controller, useForm } from "react-hook-form";
import type { z } from "zod";

import { Button } from "@loreline/ui/components/button";
import { Field, FieldError, FieldLabel } from "@loreline/ui/components/field";
import { Input } from "@loreline/ui/components/input";
import { toast } from "@loreline/ui/components/toast";

import { useForgotPassword } from "../../api";
import { forgotPasswordSchema } from "../../schemas";

export const ForgotPasswordForm = () => {
  const { mutate, isPending } = useForgotPassword();

  const form = useForm<z.infer<typeof forgotPasswordSchema>>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: {
      email: "",
    },
  });

  function onSubmit(data: z.infer<typeof forgotPasswordSchema>) {
    mutate(
      { body: data },
      {
        onSuccess: () => {
          toast.add({
            title: "Forgot password link sent",
            type: "success",
          });
        },
      },
    );
  }

  return (
    <div className="m-auto w-full max-w-md py-14">
      <p className="font-semibold text-brand-ink text-sm">Forgot password</p>
      <h2 className="mt-3 font-semibold text-5xl leading-[1.08] tracking-[-0.047em]">
        Let&apos;s get you back in.
      </h2>
      <p className="mt-3 text-muted-foreground text-sm leading-relaxed">
        Enter the email tied to your account and we&apos;ll send you a link to
        reset your password.
      </p>

      <div className="mt-9 space-y-5">
        <form
          id="form-forgot-password"
          onSubmit={form.handleSubmit(onSubmit)}
          className="space-y-5"
        >
          <Controller
            name="email"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor="form-forgot-password-email">
                  Email
                </FieldLabel>
                <Input
                  {...field}
                  id="form-forgot-password-email"
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
          <Button
            type="submit"
            size="xl"
            className="w-full"
            disabled={isPending}
          >
            {isPending ? <LoaderCircleIcon className="animate-spin" /> : null}
            Send reset link
            <ArrowRightIcon data-icon="inline-end" />
          </Button>
        </form>
      </div>

      <p className="mt-6 text-center text-muted-foreground text-sm">
        Remembered your password?{" "}
        <Link
          href="/sign-in"
          className="font-semibold text-foreground underline-offset-4 hover:underline"
        >
          Back to sign in
        </Link>
      </p>
    </div>
  );
};
