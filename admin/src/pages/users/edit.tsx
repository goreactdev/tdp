import React from "react";
import { IResourceComponentsProps } from "@refinedev/core";
import { Edit, useForm, useSelect } from "@refinedev/antd";
import { Form, Input, Select } from "antd";

export const UsersEdit: React.FC<IResourceComponentsProps> = () => {
  const { formProps, saveButtonProps, queryResult } = useForm();

  const usersData = queryResult?.data?.data;

  const { selectProps } = useSelect({
    resource: "roles",
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
          label="First Name"
          name={["first_name"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Avatar URL"
          name={["avatar_url"]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Last Name"
          name={["last_name"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Username"
          name={["username"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Friendly Address"
          name={["friendly_address"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Raw Address"
          name={["raw_address"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item label="Role" name={["role_id"]}>
          <Select
            defaultValue={usersData?.role.id}
            placeholder="Select a role"
            style={{ width: 300 }}
            {...selectProps}
          />
        </Form.Item>

        <Form.Item
          label="Job"
          name={["job"]}
        >
          <Input />
        </Form.Item>
        
        <Form.Item
          label="Bio"
          name={["bio"]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Created At"
          name={["created_at"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input readOnly disabled />
        </Form.Item>
        <Form.Item
          label="Updated At"
          name={["updated_at"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input readOnly disabled />
        </Form.Item>
        <Form.Item
          label="Awards Count"
          name={["awards_count"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input readOnly disabled />
        </Form.Item>
        <>
          {usersData?.linked_accounts?.map((item: any, index: number) => (
            <>
              <Form.Item
                key={item.id}
                label={`Linked Account ${index + 1}: Login`}
                name={["linked_accounts", index, "login"]}
                rules={[
                  {
                    required: true,
                  },
                ]}
              >
                <Input />
              </Form.Item>
              <Form.Item
                label={`Linked Account ${index + 1}: Provider`}
                name={["linked_accounts", index, "provider"]}
                rules={[
                  {
                    required: true,
                  },
                ]}
              >
                <Input />
              </Form.Item>
              <Form.Item
                label={`Linked Account ${index + 1}: Avatar URL`}
                name={["linked_accounts", index, "avatar_url"]}
              >
                <Input />
              </Form.Item>
            </>
          ))}
        </>


        <Form.Item
          label="Version"
          name={["version"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input readOnly disabled />
        </Form.Item>
      </Form>
    </Edit>
  );
};
