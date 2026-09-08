"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRightIcon } from "lucide-react";
import { Controller, useForm } from "react-hook-form";
import type { z } from "zod";

import { Button } from "@loreline/ui/components/button";
import { Field, FieldError, FieldLabel } from "@loreline/ui/components/field";
import { Input } from "@loreline/ui/components/input";
import { toast } from "@loreline/ui/components/toast";

import { resetPasswordSchema } from "../../schemas";

type ResetPasswordFormProps = {
  token: string;
};

export const ResetPasswordForm = ({ token }: ResetPasswordFormProps) => {
  const form = useForm<z.infer<typeof resetPasswordSchema>>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      token,
      password: "",
      confirmPassword: "",
    },
  });

  function onSubmit(data: z.infer<typeof resetPasswordSchema>) {
    console.log(data);
    toast.add({
      title: "Password reset successfully",
      type: "success",
    });
  }

  return (
    <div className="m-auto w-full max-w-md py-14">
      <p className="font-semibold text-brand-ink text-sm">Reset password</p>
      <h2 className="mt-3 font-semibold text-5xl leading-[1.08] tracking-[-0.047em]">
        Set a new password.
      </h2>
      <p className="mt-3 text-muted-foreground text-sm leading-relaxed">
        Make sure it&apos;s at least 6 characters and something you haven&apos;t
        used before.
      </p>

      <div className="mt-9 space-y-5">
        <form
          id="form-reset-password"
          onSubmit={form.handleSubmit(onSubmit)}
          className="space-y-5"
        >
          <Controller
            name="password"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor="form-reset-password-password">
                  New password
                </FieldLabel>
                <Input
                  {...field}
                  id="form-reset-password-password"
                  type="password"
                  aria-invalid={fieldState.invalid}
                  autoComplete="new-password"
                  placeholder="At least 6 characters"
                />
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />
          <Controller
            name="confirmPassword"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor="form-reset-password-confirm-password">
                  Confirm new password
                </FieldLabel>
                <Input
                  {...field}
                  id="form-reset-password-confirm-password"
                  type="password"
                  aria-invalid={fieldState.invalid}
                  autoComplete="new-password"
                  placeholder="Re-enter your password"
                />
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />
          <Button type="submit" size="xl" className="w-full">
            Reset password
            <ArrowRightIcon data-icon="inline-end" />
          </Button>
        </form>
      </div>
    </div>
  );
};
