import { Button, type ButtonProps } from "components/Button/Button";
import { type Ref } from "react";

export const AgentButton = ({
	ref,
	...props
}: ButtonProps & { ref?: Ref<HTMLButtonElement> }) => {
	return <Button variant="outline" ref={ref} {...props} />;
};
