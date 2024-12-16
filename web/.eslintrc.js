module.exports = {
    root: true,
    env: {
        node: true,
    },
    extends: [
        'plugin:vue/vue3-recommended',
        'eslint:recommended',
    ],
    rules: {
        'vue/multi-word-component-names': 'off', // 关闭多单词规则
    },
};
