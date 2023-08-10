import React from "react";
import { IResourceComponentsProps } from "@refinedev/core";
import { Edit, useCheckboxGroup, useForm } from "@refinedev/antd";
import { Checkbox, Form, Input } from "antd";

export const RolesEdit: React.FC<IResourceComponentsProps> = () => {
  const { formProps, saveButtonProps, queryResult } = useForm<any>();

  const { checkboxGroupProps } = useCheckboxGroup({
    resource: "permissions",
    optionLabel: "name",
    optionValue: "id",
      
  });

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Id"
          name={["id"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input readOnly disabled />
        </Form.Item>
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
        label="Tags" name="tags">
          <Checkbox.Group
          {...checkboxGroupProps} 

          />
        </Form.Item>
      </Form>
    </Edit>
  );
};
