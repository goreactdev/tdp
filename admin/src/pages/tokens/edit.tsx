import React from "react";
import { IResourceComponentsProps } from "@refinedev/core";
import { Edit, useForm, getValueFromEvent, ImageField } from "@refinedev/antd";
import { Form, Input, Upload } from "antd";

export const SBTTokensEdit: React.FC<IResourceComponentsProps> = () => {
    const { formProps, saveButtonProps, queryResult } = useForm();
    

    const nftsData = queryResult?.data?.data;
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
                    label="Raw Address"
                    name={["raw_address"]}
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <Input readOnly disabled />
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
                    <Input readOnly disabled />
                </Form.Item>
                <Form.Item
                    label="Content Uri"
                    name="content_uri"
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <Input readOnly disabled  />
                </Form.Item>
                <Form.Item
                    label="Raw Owner Address"
                    name={["raw_owner_address"]}
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <Input readOnly disabled />
                </Form.Item>
                <Form.Item
                    label="Friendly Owner Address"
                    name={["friendly_owner_address"]}
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
                    <Input readOnly disabled />
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
                    <Input readOnly disabled />
                </Form.Item>
                <Form.Item label="Image">
                    <Form.Item
                        name="image"

                    >
                         <Input readOnly disabled />
                    </Form.Item>
                    <ImageField style={{maxWidth: 400}} value={nftsData?.image} />
                </Form.Item>
                <Form.Item
                    label="Rating points"
                    name={["weight"]}
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="Index"
                    name={["index"]}
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <Input readOnly disabled />
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