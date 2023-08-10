import React from "react";
import { IResourceComponentsProps } from "@refinedev/core";
import { Create, useCheckboxGroup, useForm } from "@refinedev/antd";
import { Checkbox, Form, Input } from "antd";

export const RolesCreate: React.FC<IResourceComponentsProps> = () => {
  const { formProps, saveButtonProps } = useForm();
  const { checkboxGroupProps } = useCheckboxGroup({
    resource: "permissions",
    optionLabel: "name",
    optionValue: "id",    
  });

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Name"
          name={["name"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Description"
          name={["description"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          rules={[
            {
              required: true,
            },
          ]}
          label="Permissions"
          name="permissions"
        >
          <Checkbox.Group {...checkboxGroupProps} />
        </Form.Item>
      </Form>
    </Create>
  );
};
