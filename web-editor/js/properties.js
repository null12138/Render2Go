// 属性面板系统 - properties.js

// 属性控件基类
class PropertyControl extends EventEmitter {
    constructor(property, label, value) {
        super();
        this.property = property;
        this.label = label;
        this.value = value;
        this.element = null;
        this.createElement();
    }

    createElement() {
        this.element = Utils.createElement('div', 'property-control');
        this.element.innerHTML = `
            <label class="property-label">${this.label}</label>
            <div class="property-input-container">
                ${this.createInput()}
            </div>
        `;
        this.bindEvents();
    }

    createInput() {
        return '<input type="text" class="property-input">';
    }

    bindEvents() {
        const input = this.element.querySelector('.property-input');
        if (input) {
            input.value = this.value;
            input.addEventListener('change', () => {
                this.setValue(input.value);
            });
        }
    }

    setValue(value) {
        this.value = value;
        this.emit('change', this.property, value);

        const input = this.element.querySelector('.property-input');
        if (input && input.value !== value) {
            input.value = value;
        }
    }

    getValue() {
        return this.value;
    }
}

// 数字控件
class NumberControl extends PropertyControl {
    constructor(property, label, value, min = null, max = null, step = 1) {
        super(property, label, value);
        this.min = min;
        this.max = max;
        this.step = step;
    }

    createInput() {
        const attrs = [];
        if (this.min !== null) attrs.push(`min="${this.min}"`);
        if (this.max !== null) attrs.push(`max="${this.max}"`);
        attrs.push(`step="${this.step}"`);

        return `<input type="number" class="property-input" ${attrs.join(' ')}>`;
    }

    setValue(value) {
        const numValue = parseFloat(value);
        if (!isNaN(numValue)) {
            let clampedValue = numValue;
            if (this.min !== null) clampedValue = Math.max(this.min, clampedValue);
            if (this.max !== null) clampedValue = Math.min(this.max, clampedValue);
            super.setValue(clampedValue);
        }
    }
}

// 颜色控件
class ColorControl extends PropertyControl {
    createInput() {
        return `
            <div class="color-input-group">
                <input type="color" class="color-picker" value="${this.value}">
                <input type="text" class="property-input color-text" value="${this.value}">
            </div>
        `;
    }

    bindEvents() {
        const colorPicker = this.element.querySelector('.color-picker');
        const textInput = this.element.querySelector('.color-text');

        if (colorPicker) {
            colorPicker.addEventListener('input', () => {
                this.setValue(colorPicker.value);
                textInput.value = colorPicker.value;
            });
        }

        if (textInput) {
            textInput.addEventListener('change', () => {
                this.setValue(textInput.value);
                colorPicker.value = textInput.value;
            });
        }
    }
}

// 选择控件
class SelectControl extends PropertyControl {
    constructor(property, label, value, options) {
        super(property, label, value);
        this.options = options;
    }

    createInput() {
        const optionsHtml = this.options.map(option => {
            const selected = option.value === this.value ? 'selected' : '';
            return `<option value="${option.value}" ${selected}>${option.label}</option>`;
        }).join('');

        return `<select class="property-input">${optionsHtml}</select>`;
    }

    bindEvents() {
        const select = this.element.querySelector('select');
        if (select) {
            select.addEventListener('change', () => {
                this.setValue(select.value);
            });
        }
    }
}

// 复选框控件
class CheckboxControl extends PropertyControl {
    createInput() {
        const checked = this.value ? 'checked' : '';
        return `<input type="checkbox" class="property-checkbox" ${checked}>`;
    }

    bindEvents() {
        const checkbox = this.element.querySelector('.property-checkbox');
        if (checkbox) {
            checkbox.addEventListener('change', () => {
                this.setValue(checkbox.checked);
            });
        }
    }
}

// 文本控件
class TextControl extends PropertyControl {
    constructor(property, label, value, multiline = false) {
        super(property, label, value);
        this.multiline = multiline;
    }

    createInput() {
        if (this.multiline) {
            return `<textarea class="property-input property-textarea" rows="3">${this.value}</textarea>`;
        }
        return `<input type="text" class="property-input" value="${this.value}">`;
    }

    bindEvents() {
        const input = this.element.querySelector('.property-input');
        if (input) {
            input.addEventListener('input', () => {
                this.setValue(input.value);
            });
        }
    }
}

// 滑块控件
class SliderControl extends PropertyControl {
    constructor(property, label, value, min = 0, max = 100, step = 1) {
        super(property, label, value);
        this.min = min;
        this.max = max;
        this.step = step;
    }

    createInput() {
        return `
            <div class="slider-input-group">
                <input type="range" class="property-slider" 
                       min="${this.min}" max="${this.max}" step="${this.step}" value="${this.value}">
                <input type="number" class="property-input slider-number" 
                       min="${this.min}" max="${this.max}" step="${this.step}" value="${this.value}">
            </div>
        `;
    }

    bindEvents() {
        const slider = this.element.querySelector('.property-slider');
        const numberInput = this.element.querySelector('.slider-number');

        const updateValue = (value) => {
            this.setValue(parseFloat(value));
            slider.value = this.value;
            numberInput.value = this.value;
        };

        if (slider) {
            slider.addEventListener('input', () => updateValue(slider.value));
        }

        if (numberInput) {
            numberInput.addEventListener('change', () => updateValue(numberInput.value));
        }
    }
}

// 属性组
class PropertyGroup {
    constructor(title, collapsed = false) {
        this.title = title;
        this.collapsed = collapsed;
        this.controls = [];
        this.element = null;
        this.createElement();
    }

    createElement() {
        this.element = Utils.createElement('div', 'property-group');
        this.element.innerHTML = `
            <div class="property-group-header ${this.collapsed ? 'collapsed' : ''}">
                <i class="fas fa-chevron-down group-toggle"></i>
                <span class="group-title">${this.title}</span>
            </div>
            <div class="property-group-content ${this.collapsed ? 'hidden' : ''}"></div>
        `;

        const header = this.element.querySelector('.property-group-header');
        header.addEventListener('click', () => this.toggle());
    }

    addControl(control) {
        this.controls.push(control);
        const content = this.element.querySelector('.property-group-content');
        content.appendChild(control.element);
        return control;
    }

    removeControl(control) {
        const index = this.controls.indexOf(control);
        if (index >= 0) {
            this.controls.splice(index, 1);
            control.element.remove();
        }
    }

    toggle() {
        this.collapsed = !this.collapsed;
        const header = this.element.querySelector('.property-group-header');
        const content = this.element.querySelector('.property-group-content');

        header.classList.toggle('collapsed', this.collapsed);
        content.classList.toggle('hidden', this.collapsed);
    }

    clear() {
        this.controls.forEach(control => control.element.remove());
        this.controls = [];
    }
}

// 属性面板
class PropertiesPanel extends EventEmitter {
    constructor() {
        super();
        this.container = null;
        this.groups = [];
        this.selectedObjects = [];
        this.isUpdating = false;

        this.init();
    }

    init() {
        this.container = document.getElementById('propertiesContent');
        if (this.container) {
            this.setupDefaultGroups();
            // 延迟绑定事件，确保app存在
            setTimeout(() => this.bindObjectManagerEvents(), 100);
        }
    }

    setupDefaultGroups() {
        this.transformGroup = new PropertyGroup('变换', false);
        this.appearanceGroup = new PropertyGroup('外观', false);
        this.textGroup = new PropertyGroup('文本', true);
        this.animationGroup = new PropertyGroup('动画', true);

        this.groups = [
            this.transformGroup,
            this.appearanceGroup,
            this.textGroup,
            this.animationGroup
        ];

        this.updateDisplay();
    }

    bindObjectManagerEvents() {
        if (window.app && window.app.objectManager) {
            window.app.objectManager.on('selectionChanged', (selectedObjects) => {
                this.setSelectedObjects(selectedObjects);
            });

            window.app.objectManager.on('objectUpdated', (object) => {
                if (this.selectedObjects.includes(object)) {
                    this.updateProperties();
                }
            });
        }
    }

    setSelectedObjects(objects) {
        this.selectedObjects = objects;
        this.updateProperties();
    }

    updateProperties() {
        if (this.isUpdating) return;

        this.clearControls();

        if (this.selectedObjects.length === 0) {
            this.showNoSelection();
            return;
        }

        if (this.selectedObjects.length === 1) {
            this.showSingleObjectProperties(this.selectedObjects[0]);
        } else {
            this.showMultiObjectProperties(this.selectedObjects);
        }

        this.updateDisplay();
    }

    clearControls() {
        this.groups.forEach(group => group.clear());
    }

    showNoSelection() {
        const message = Utils.createElement('div', 'no-selection-message');
        message.textContent = '未选择对象';
        this.container.innerHTML = '';
        this.container.appendChild(message);
    }

    showSingleObjectProperties(object) {
        // 变换属性
        this.addTransformControls(object);

        // 外观属性
        this.addAppearanceControls(object);

        // 文本属性（如果是文本对象）
        if (object.type === 'text') {
            this.addTextControls(object);
        }

        // 动画属性
        this.addAnimationControls(object);
    }

    showMultiObjectProperties(objects) {
        // 只显示通用属性
        this.addCommonTransformControls(objects);
        this.addCommonAppearanceControls(objects);
    }

    addTransformControls(object) {
        const xControl = new NumberControl('x', 'X 位置', object.x, null, null, 1);
        const yControl = new NumberControl('y', 'Y 位置', object.y, null, null, 1);
        const widthControl = new NumberControl('width', '宽度', object.width, 1, null, 1);
        const heightControl = new NumberControl('height', '高度', object.height, 1, null, 1);
        const rotationControl = new NumberControl('rotation', '旋转', object.rotation || 0, -360, 360, 1);

        [xControl, yControl, widthControl, heightControl, rotationControl].forEach(control => {
            control.on('change', (property, value) => this.updateObjectProperty(object, property, value));
            this.transformGroup.addControl(control);
        });
    }

    addAppearanceControls(object) {
        const opacityControl = new SliderControl('opacity', '不透明度', object.opacity || 1, 0, 1, 0.01);
        const visibleControl = new CheckboxControl('visible', '可见', object.visible !== false);
        const lockedControl = new CheckboxControl('locked', '锁定', object.locked === true);

        [opacityControl, visibleControl, lockedControl].forEach(control => {
            control.on('change', (property, value) => this.updateObjectProperty(object, property, value));
            this.appearanceGroup.addControl(control);
        });

        // 颜色属性（根据对象类型）
        if (object.fillColor !== undefined) {
            const fillColorControl = new ColorControl('fillColor', '填充颜色', object.fillColor);
            fillColorControl.on('change', (property, value) => this.updateObjectProperty(object, property, value));
            this.appearanceGroup.addControl(fillColorControl);
        }

        if (object.strokeColor !== undefined) {
            const strokeColorControl = new ColorControl('strokeColor', '描边颜色', object.strokeColor);
            const strokeWidthControl = new NumberControl('strokeWidth', '描边宽度', object.strokeWidth || 1, 0, null, 1);

            strokeColorControl.on('change', (property, value) => this.updateObjectProperty(object, property, value));
            strokeWidthControl.on('change', (property, value) => this.updateObjectProperty(object, property, value));

            this.appearanceGroup.addControl(strokeColorControl);
            this.appearanceGroup.addControl(strokeWidthControl);
        }
    }

    addTextControls(object) {
        const textControl = new TextControl('text', '文本内容', object.text || '', true);
        const fontSizeControl = new NumberControl('fontSize', '字体大小', object.fontSize || 16, 8, 200, 1);
        const fontFamilyControl = new SelectControl('fontFamily', '字体', object.fontFamily || 'Arial', [
            { value: 'Arial', label: 'Arial' },
            { value: 'Helvetica', label: 'Helvetica' },
            { value: 'Times New Roman', label: 'Times New Roman' },
            { value: 'Courier New', label: 'Courier New' },
            { value: 'Microsoft YaHei', label: '微软雅黑' },
            { value: 'SimSun', label: '宋体' },
            { value: 'SimHei', label: '黑体' }
        ]);

        const alignControl = new SelectControl('textAlign', '对齐方式', object.textAlign || 'left', [
            { value: 'left', label: '左对齐' },
            { value: 'center', label: '居中' },
            { value: 'right', label: '右对齐' }
        ]);

        [textControl, fontSizeControl, fontFamilyControl, alignControl].forEach(control => {
            control.on('change', (property, value) => this.updateObjectProperty(object, property, value));
            this.textGroup.addControl(control);
        });

        // 展开文本组
        this.textGroup.collapsed = false;
    }

    addAnimationControls(object) {
        // 添加关键帧按钮
        const addKeyframeBtn = Utils.createElement('button', 'add-keyframe-btn');
        addKeyframeBtn.textContent = '添加关键帧';
        addKeyframeBtn.addEventListener('click', () => {
            if (app.timeline) {
                app.timeline.autoKeyframe(object.id);
            }
        });

        const buttonContainer = Utils.createElement('div', 'property-control');
        buttonContainer.appendChild(addKeyframeBtn);
        this.animationGroup.addControl({ element: buttonContainer });
    }

    addCommonTransformControls(objects) {
        // 获取公共值
        const commonX = this.getCommonValue(objects, 'x');
        const commonY = this.getCommonValue(objects, 'y');
        const commonWidth = this.getCommonValue(objects, 'width');
        const commonHeight = this.getCommonValue(objects, 'height');

        if (commonX !== null) {
            const xControl = new NumberControl('x', 'X 位置', commonX);
            xControl.on('change', (property, value) => this.updateMultipleObjects(objects, property, value));
            this.transformGroup.addControl(xControl);
        }

        if (commonY !== null) {
            const yControl = new NumberControl('y', 'Y 位置', commonY);
            yControl.on('change', (property, value) => this.updateMultipleObjects(objects, property, value));
            this.transformGroup.addControl(yControl);
        }

        // 对齐按钮
        this.addAlignmentControls(objects);
    }

    addCommonAppearanceControls(objects) {
        const commonOpacity = this.getCommonValue(objects, 'opacity');

        if (commonOpacity !== null) {
            const opacityControl = new SliderControl('opacity', '不透明度', commonOpacity, 0, 1, 0.01);
            opacityControl.on('change', (property, value) => this.updateMultipleObjects(objects, property, value));
            this.appearanceGroup.addControl(opacityControl);
        }
    }

    addAlignmentControls(objects) {
        const alignmentContainer = Utils.createElement('div', 'alignment-controls');
        alignmentContainer.innerHTML = `
            <div class="property-label">对齐</div>
            <div class="alignment-buttons">
                <button class="align-btn" data-align="left">左对齐</button>
                <button class="align-btn" data-align="center">居中</button>
                <button class="align-btn" data-align="right">右对齐</button>
                <button class="align-btn" data-align="top">顶部对齐</button>
                <button class="align-btn" data-align="middle">垂直居中</button>
                <button class="align-btn" data-align="bottom">底部对齐</button>
            </div>
        `;

        alignmentContainer.querySelectorAll('.align-btn').forEach(btn => {
            btn.addEventListener('click', () => {
                app.objectManager.alignObjects(objects.map(obj => obj.id), btn.dataset.align);
            });
        });

        this.transformGroup.addControl({ element: alignmentContainer });
    }

    getCommonValue(objects, property) {
        if (objects.length === 0) return null;

        const firstValue = objects[0][property];
        const allSame = objects.every(obj => obj[property] === firstValue);

        return allSame ? firstValue : null;
    }

    updateObjectProperty(object, property, value) {
        this.isUpdating = true;

        // 记录撤销状态
        app.undoManager.execute(new Command(
            () => this.applyObjectProperty(object, property, value),
            () => this.applyObjectProperty(object, property, object[property])
        ));

        this.isUpdating = false;
    }

    updateMultipleObjects(objects, property, value) {
        this.isUpdating = true;

        const oldValues = objects.map(obj => obj[property]);

        app.undoManager.execute(new Command(
            () => objects.forEach(obj => this.applyObjectProperty(obj, property, value)),
            () => objects.forEach((obj, i) => this.applyObjectProperty(obj, property, oldValues[i]))
        ));

        this.isUpdating = false;
    }

    applyObjectProperty(object, property, value) {
        switch (property) {
            case 'x':
            case 'y':
                object.setPosition(
                    property === 'x' ? value : object.x,
                    property === 'y' ? value : object.y
                );
                break;
            case 'width':
            case 'height':
                object.setSize(
                    property === 'width' ? value : object.width,
                    property === 'height' ? value : object.height
                );
                break;
            case 'rotation':
                object.setRotation(value);
                break;
            case 'opacity':
                object.setOpacity(value);
                break;
            case 'visible':
                object.setVisible(value);
                break;
            case 'locked':
                object.setLocked(value);
                break;
            case 'fillColor':
                object.setFillColor(value);
                break;
            case 'strokeColor':
                object.setStrokeColor(value);
                break;
            case 'strokeWidth':
                object.setStrokeWidth(value);
                break;
            case 'text':
                if (object.setText) object.setText(value);
                break;
            case 'fontSize':
                if (object.setFontSize) object.setFontSize(value);
                break;
            case 'fontFamily':
                if (object.setFontFamily) object.setFontFamily(value);
                break;
            case 'textAlign':
                if (object.setTextAlign) object.setTextAlign(value);
                break;
            default:
                object[property] = value;
                object.updateElement();
        }
    }

    updateDisplay() {
        this.container.innerHTML = '';
        this.groups.forEach(group => {
            if (group.controls.length > 0) {
                this.container.appendChild(group.element);
            }
        });
    }

    // 外部API
    refresh() {
        this.updateProperties();
    }

    showGroup(groupName) {
        const group = this.groups.find(g => g.title === groupName);
        if (group) {
            group.collapsed = false;
            group.toggle();
        }
    }

    hideGroup(groupName) {
        const group = this.groups.find(g => g.title === groupName);
        if (group) {
            group.collapsed = true;
            group.toggle();
        }
    }
}

// 导出
window.PropertiesPanel = PropertiesPanel;
window.PropertyControl = PropertyControl;
window.NumberControl = NumberControl;
window.ColorControl = ColorControl;
window.SelectControl = SelectControl;
window.CheckboxControl = CheckboxControl;
window.TextControl = TextControl;
window.SliderControl = SliderControl;
window.PropertyGroup = PropertyGroup;
