/**
 * Copied from shadc/ui on 03/05/2025
 * @see {@link https://ui.shadcn.com/docs/components/scroll-area}
 */
import * as ScrollAreaPrimitive from "@radix-ui/react-scroll-area";
import * as React from "react";
import { cn } from "utils/cn";

export const ScrollArea: React.FC<
	React.ComponentPropsWithRef<typeof ScrollAreaPrimitive.Root>
> = ({ className, children, ref, ...props }) => {
	return (
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
};

export const ScrollBar: React.FC<
	React.ComponentPropsWithRef<typeof ScrollAreaPrimitive.ScrollAreaScrollbar>
> = ({ className, orientation = "vertical", ref, ...props }) => {
	return (
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
};
