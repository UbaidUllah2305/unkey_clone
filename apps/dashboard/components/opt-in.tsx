"use client";

import { trpc } from "@/lib/trpc/client";
import { PostHogEvent } from "@/providers/PostHogProvider";
import type { Workspace } from "@unkey/db";
import { Button, Empty, toast } from "@unkey/ui";
import Link from "next/link";
import { useRouter } from "next/navigation";
import type React from "react";

type Props = {
  title: string;
  description: string;
  feature: keyof Workspace["betaFeatures"];
};

export const OptIn: React.FC<Props> = ({ title, description, feature }) => {
  const router = useRouter();
  const optIn = trpc.workspace.optIntoBeta.useMutation({
    onMutate() {
      PostHogEvent({
        name: "self-serve-opt-in",
        properties: { feature },
      });
    },
    onSuccess() {
      PostHogEvent({
        name: "self-serve-opt-in",
        properties: { feature },
      });
      toast.success("Successfully opted in");
      router.refresh();
    },
    onError(err) {
      toast.error(err.message);
    },
  });
  return (
    <Empty>
      <Empty.Icon />
      <Empty.Title>{title}</Empty.Title>
      <Empty.Description>{description}</Empty.Description>

      <div className="flex items-center gap-4">
        <Link href={`mailto:support@unkey.dev?subject=Beta Access: ${feature}`}>
          <Button>Get in touch</Button>
        </Link>

        <Button
          variant="primary"
          onClick={() => optIn.mutate({ feature })}
          loading={optIn.isLoading}
        >
          Opt In
        </Button>
      </div>
    </Empty>
  );
};
