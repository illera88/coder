import { css } from "@emotion/css";
import { useTheme } from "@emotion/react";
import IconButton, { type IconButtonProps } from "@mui/material/IconButton";
import { Avatar, type AvatarProps } from "components/Avatar/Avatar";
import { Button, type ButtonProps } from "components/Button/Button";
import {
	cloneElement,
	type FC,
	type HTMLAttributes,
	type ReactElement,
} from "react";
import { cn } from "utils/cn";

export const Topbar: FC<HTMLAttributes<HTMLElement>> = (props) => {
	const theme = useTheme();

	return (
		<header
			{...props}
			css={{
				minHeight: 48,
				borderBottom: `1px solid ${theme.palette.divider}`,
				display: "flex",
				alignItems: "center",
				fontSize: 13,
				lineHeight: "1.2",
			}}
		/>
	);
};

export const TopbarIconButton = ({
	ref,
	...props
}: IconButtonProps & { ref?: Ref<HTMLButtonElement> }) => {
	return (
		<IconButton
			ref={ref}
			{...props}
			size="small"
			css={{
				padding: 0,
				borderRadius: 0,
				height: 48,
				width: 48,

				"& svg": {
					fontSize: 20,
				},
			}}
		/>
	);
};

export const TopbarButton = ({
	ref,
	...props
}: ButtonProps & { ref?: Ref<HTMLButtonElement> }) => {
	return <Button ref={ref} variant="outline" size="sm" {...props} />;
};

export const TopbarData: FC<HTMLAttributes<HTMLDivElement>> = (props) => {
	return (
		<div
			{...props}
			css={{
				display: "flex",
				gap: 8,
				alignItems: "center",
				justifyContent: "center",
			}}
		/>
	);
};

export const TopbarDivider: FC<HTMLAttributes<HTMLSpanElement>> = (props) => {
	const theme = useTheme();
	return (
		<span {...props} css={{ color: theme.palette.divider }}>
			/
		</span>
	);
};

export const TopbarAvatar: FC<AvatarProps> = (props) => {
	return <Avatar {...props} variant="icon" size="md" />;
};

type TopbarIconProps = HTMLAttributes<HTMLOrSVGElement> & {
	ref?: React.Ref<HTMLOrSVGElement>;
};

export const TopbarIcon = ({
	children,
	className,
	ref,
	...restProps
}: TopbarIconProps) => {
	const theme = useTheme();

	return cloneElement(
		children as ReactElement<
			HTMLAttributes<HTMLOrSVGElement> & {
				ref?: React.Ref<HTMLOrSVGElement>;
			}
		>,
		{
			...restProps,
			ref,
			className: cn([
				css({ fontSize: 16, color: theme.palette.text.disabled }),
				"size-icon-sm",
			]),
		},
	);
};
