/**
 * Copied from shadc/ui on 11/13/2024
 * @see {@link https://ui.shadcn.com/docs/components/label}
 */
import * as LabelPrimitive from "@radix-ui/react-label";
import { cva, type VariantProps } from "class-variance-authority";
import * as React from "react";
import { cn } from "utils/cn";

const labelVariants = cva(
	"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
);

type LabelProps = React.ComponentPropsWithRef<typeof LabelPrimitive.Root> &
	VariantProps<typeof labelVariants>;

export const Label: FC<LabelProps> = ({ className, ...props }) => (
	<LabelPrimitive.Root className={cn(labelVariants(), className)} {...props} />
);
