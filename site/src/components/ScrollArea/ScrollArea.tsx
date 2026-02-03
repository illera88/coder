/**
 * Copied from shadc/ui on 03/05/2025
 * @see {@link https://ui.shadcn.com/docs/components/scroll-area}
 */
import * as ScrollAreaPrimitive from "@radix-ui/react-scroll-area";
import * as React from "react";
import { cn } from "utils/cn";

export const ScrollArea = ({
	className,
	children,
	ref,
	...props
}: React.ComponentPropsWithoutRef<typeof ScrollAreaPrimitive.Root> & {
	ref?: React.Ref<React.ComponentRef<typeof ScrollAreaPrimitive.Root>>;
}) => (
	<ScrollAreaPrimitive.Root
		ref={ref}
		className={cn("relative overflow-hidden", className)}
		{...props}
	>
		<ScrollAreaPrimitive.Viewport className="h-full w-full rounded-[inherit]">
			{children}
		</ScrollAreaPrimitive.Viewport>
		<ScrollBar className="z-10" />
		<ScrollAreaPrimitive.Corner />
	</ScrollAreaPrimitive.Root>
);
ScrollArea.displayName = ScrollAreaPrimitive.Root.displayName;

export const ScrollBar = ({
	className,
	orientation = "vertical",
	ref,
	...props
}: React.ComponentPropsWithoutRef<
	typeof ScrollAreaPrimitive.ScrollAreaScrollbar
> & {
	ref?: React.Ref<
		React.ComponentRef<typeof ScrollAreaPrimitive.ScrollAreaScrollbar>
	>;
}) => (
	<ScrollAreaPrimitive.ScrollAreaScrollbar
		ref={ref}
		orientation={orientation}
		className={cn(
			"border-0 border-solid border-border flex touch-none select-none transition-colors",
			orientation === "vertical" &&
				"h-full w-2.5 border-l border-l-transparent p-[1px]",
			orientation === "horizontal" &&
				"h-2.5 flex-col border-t border-t-transparent p-[1px]",
			className,
		)}
		{...props}
	>
		<ScrollAreaPrimitive.ScrollAreaThumb className="relative flex-1 rounded-full bg-surface-quaternary" />
	</ScrollAreaPrimitive.ScrollAreaScrollbar>
);
