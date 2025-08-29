// 工具函数 - utils.js

// 全局工具函数
const Utils = {
    // 生成唯一ID
    generateId() {
        return 'obj_' + Math.random().toString(36).substr(2, 9);
    },

    // 颜色转换
    hexToRgb(hex) {
        const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
        return result ? {
            r: parseInt(result[1], 16),
            g: parseInt(result[2], 16),
            b: parseInt(result[3], 16)
        } : null;
    },

    rgbToHex(r, g, b) {
        return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
    },

    // 数学工具
    clamp(value, min, max) {
        return Math.min(Math.max(value, min), max);
    },

    lerp(start, end, t) {
        return start + (end - start) * t;
    },

    distance(x1, y1, x2, y2) {
        return Math.sqrt(Math.pow(x2 - x1, 2) + Math.pow(y2 - y1, 2));
    },

    // 角度转换
    degToRad(degrees) {
        return degrees * (Math.PI / 180);
    },

    radToDeg(radians) {
        return radians * (180 / Math.PI);
    },

    // DOM操作
    createElement(tag, className, parent) {
        const element = document.createElement(tag);
        if (className) element.className = className;
        if (parent) parent.appendChild(element);
        return element;
    },

    // 事件处理
    on(element, event, handler) {
        element.addEventListener(event, handler);
        return () => element.removeEventListener(event, handler);
    },

    // 防抖函数
    debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    },

    // 节流函数
    throttle(func, limit) {
        let inThrottle;
        return function () {
            const args = arguments;
            const context = this;
            if (!inThrottle) {
                func.apply(context, args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        }
    },

    // 深拷贝
    deepClone(obj) {
        if (obj === null || typeof obj !== "object") return obj;
        if (obj instanceof Date) return new Date(obj.getTime());
        if (obj instanceof Array) return obj.map(item => this.deepClone(item));
        if (typeof obj === "object") {
            const clonedObj = {};
            for (const key in obj) {
                if (obj.hasOwnProperty(key)) {
                    clonedObj[key] = this.deepClone(obj[key]);
                }
            }
            return clonedObj;
        }
    },

    // 格式化时间
    formatTime(seconds) {
        const mins = Math.floor(seconds / 60);
        const secs = Math.floor(seconds % 60);
        return `${mins}:${secs.toString().padStart(2, '0')}`;
    },

    // 文件操作
    downloadFile(content, filename, type = 'text/plain') {
        const blob = new Blob([content], { type });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    },

    // 读取文件
    readFile(file) {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = e => resolve(e.target.result);
            reader.onerror = reject;
            reader.readAsText(file);
        });
    },

    // 坐标转换
    screenToCanvas(x, y, canvas) {
        const rect = canvas.getBoundingClientRect();
        return {
            x: x - rect.left,
            y: y - rect.top
        };
    },

    canvasToScreen(x, y, canvas) {
        const rect = canvas.getBoundingClientRect();
        return {
            x: x + rect.left,
            y: y + rect.top
        };
    },

    // 检测点是否在矩形内
    pointInRect(px, py, x, y, width, height) {
        return px >= x && px <= x + width && py >= y && py <= y + height;
    },

    // 检测点是否在圆内
    pointInCircle(px, py, cx, cy, radius) {
        return this.distance(px, py, cx, cy) <= radius;
    },

    // 格式化数字
    formatNumber(num, decimals = 2) {
        return Number(Math.round(num + 'e' + decimals) + 'e-' + decimals);
    },

    // 显示通知
    showNotification(message, type = 'info', duration = 3000) {
        const notification = this.createElement('div', `notification ${type}`);
        notification.textContent = message;
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 12px 20px;
            border-radius: 6px;
            color: white;
            font-weight: 500;
            z-index: 10000;
            transform: translateX(100%);
            transition: transform 0.3s ease;
        `;

        switch (type) {
            case 'success':
                notification.style.background = '#059669';
                break;
            case 'warning':
                notification.style.background = '#d97706';
                break;
            case 'error':
                notification.style.background = '#dc2626';
                break;
            default:
                notification.style.background = '#2563eb';
        }

        document.body.appendChild(notification);

        // 动画显示
        setTimeout(() => {
            notification.style.transform = 'translateX(0)';
        }, 10);

        // 自动隐藏
        setTimeout(() => {
            notification.style.transform = 'translateX(100%)';
            setTimeout(() => {
                document.body.removeChild(notification);
            }, 300);
        }, duration);
    }
};

// 事件发布订阅系统
class EventEmitter {
    constructor() {
        this.events = {};
    }

    on(event, callback) {
        if (!this.events[event]) {
            this.events[event] = [];
        }
        this.events[event].push(callback);

        // 返回取消订阅函数
        return () => {
            this.events[event] = this.events[event].filter(cb => cb !== callback);
        };
    }

    emit(event, ...args) {
        if (this.events[event]) {
            this.events[event].forEach(callback => {
                try {
                    callback(...args);
                } catch (error) {
                    console.error('Event callback error:', error);
                }
            });
        }
    }

    off(event, callback) {
        if (this.events[event]) {
            this.events[event] = this.events[event].filter(cb => cb !== callback);
        }
    }

    once(event, callback) {
        const onceCallback = (...args) => {
            callback(...args);
            this.off(event, onceCallback);
        };
        this.on(event, onceCallback);
    }
}

// 撤销重做系统
class UndoRedoManager extends EventEmitter {
    constructor(maxHistorySize = 50) {
        super();
        this.history = [];
        this.currentIndex = -1;
        this.maxHistorySize = maxHistorySize;
    }

    execute(command) {
        // 删除当前位置之后的所有历史
        this.history = this.history.slice(0, this.currentIndex + 1);

        // 添加新命令
        this.history.push(command);
        this.currentIndex++;

        // 限制历史大小
        if (this.history.length > this.maxHistorySize) {
            this.history.shift();
            this.currentIndex--;
        }

        // 执行命令
        command.execute();
        this.emit('executed', command);
    }

    undo() {
        if (this.canUndo()) {
            const command = this.history[this.currentIndex];
            command.undo();
            this.currentIndex--;
            this.emit('undone', command);
            return true;
        }
        return false;
    }

    redo() {
        if (this.canRedo()) {
            this.currentIndex++;
            const command = this.history[this.currentIndex];
            command.execute();
            this.emit('redone', command);
            return true;
        }
        return false;
    }

    canUndo() {
        return this.currentIndex >= 0;
    }

    canRedo() {
        return this.currentIndex < this.history.length - 1;
    }

    clear() {
        this.history = [];
        this.currentIndex = -1;
    }
}

// 简单的命令类
class Command {
    constructor(executeFunc, undoFunc) {
        this.executeFunc = executeFunc;
        this.undoFunc = undoFunc;
    }

    execute() {
        this.executeFunc();
    }

    undo() {
        this.undoFunc();
    }
}

// 键盘快捷键管理器
class KeyboardManager {
    constructor() {
        this.shortcuts = new Map();
        this.init();
    }

    init() {
        document.addEventListener('keydown', (e) => {
            const key = this.getKeyString(e);
            const callback = this.shortcuts.get(key);
            if (callback) {
                e.preventDefault();
                callback(e);
            }
        });
    }

    getKeyString(e) {
        const parts = [];
        if (e.ctrlKey) parts.push('Ctrl');
        if (e.altKey) parts.push('Alt');
        if (e.shiftKey) parts.push('Shift');
        if (e.metaKey) parts.push('Meta');
        parts.push(e.key);
        return parts.join('+');
    }

    register(keyString, callback) {
        this.shortcuts.set(keyString, callback);
    }

    unregister(keyString) {
        this.shortcuts.delete(keyString);
    }
}

// 导出全局对象
window.Utils = Utils;
window.EventEmitter = EventEmitter;
window.UndoRedoManager = UndoRedoManager;
window.Command = Command;
window.KeyboardManager = KeyboardManager;
